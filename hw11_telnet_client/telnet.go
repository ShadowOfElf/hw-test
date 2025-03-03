package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

var (
	errorClose     = errors.New("attempt to close a non-existent connection")
	errorConnect   = errors.New("attempt to connect fail")
	errorWrite     = errors.New("error writing to socket")
	errorReadInput = errors.New("error reading from input")
)

type client struct {
	conn    net.Conn
	in      io.Reader
	out     io.Writer
	addr    string
	timeout time.Duration
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	telnet := client{
		in:      in,
		out:     out,
		addr:    address,
		timeout: timeout,
	}

	return &telnet
}

func (c *client) Connect() error {
	if c.conn != nil {
		return nil
	}
	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return err
	}

	c.conn = conn
	if c.conn == nil {
		return errorConnect
	}

	_, err = fmt.Fprintf(os.Stderr, "connected to host: %s\n", c.addr)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Send() error {
	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		line := scanner.Text() + "\n"
		_, err := c.conn.Write([]byte(fmt.Sprintf("%s", line)))
		if err != nil {
			return errorWrite
		}
	}
	if err := scanner.Err(); err != nil && errors.Is(err, io.EOF) {
		return errorReadInput
	}
	return nil
}

func (c *client) Receive() error {
	reader := bufio.NewReader(c.conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		_, err = fmt.Fprint(c.out, line)
		if err != nil {
			return err
		}
	}
}

func (c *client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	_, err := fmt.Fprintf(os.Stderr, "close connected to host: %s\n", c.addr)
	if err != nil {
		return err
	}
	return errorClose
}

func TelnetWork(c TelnetClient) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err := c.Connect()
	if err != nil {
		log.Fatal("Connect to host error")
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := c.Send(); err != nil {
			return // TODO: добавить обработку ошибок
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.Receive(); err != nil {
			return // TODO: добавить обработку ошибок
		}
	}()

	ch := make(chan struct{})
	go func() {
		for {
			select {
			case <-ch:
				return
			case <-ctx.Done():
				_ = c.Close()
			}
		}
	}()

	wg.Wait()
	close(ch)
	_, err = fmt.Fprint(os.Stderr, "telnet work ended")

	if err != nil {
		_ = c.Close()
		return err
	}
	_ = c.Close()
	return nil
}
