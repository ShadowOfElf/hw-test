package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

var (
	errorClose         = errors.New("attempt to close a non-existent connection")
	errorDoubleConnect = errors.New("connection already established ")
)

type clientT struct {
	conn    net.Conn
	in      io.Reader
	out     io.Writer
	addr    string
	timeout time.Duration
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	telnet := clientT{
		in:      in,
		out:     out,
		addr:    address,
		timeout: timeout,
	}

	return &telnet
}

func (c *clientT) Connect() error {
	if c.conn != nil {
		return errorDoubleConnect
	}
	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return err
	}

	c.conn = conn

	_, err = fmt.Fprintf(os.Stderr, "connected to host: %s\n", c.addr)
	if err != nil {
		return err
	}

	return nil
}

func (c *clientT) Send() error {
	scanner := bufio.NewScanner(c.in)

	for scanner.Scan() {
		select {
		case <-done:
			_, _ = fmt.Fprintf(os.Stderr, "received stop signal, exiting... \n")
			return nil
		default:
			line := scanner.Text() + "\n"
			_, err := c.conn.Write([]byte(line))
			if err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil && errors.Is(err, io.EOF) {
		_, _ = fmt.Fprintf(os.Stderr, "Received EOF (Ctrl+D)...\n")
		return nil
	}
	return nil
}

func (c *clientT) Receive() error {
	reader := bufio.NewReader(c.conn)

	for {
		line, err := reader.ReadString('\n')
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

func (c *clientT) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return errorClose
}
