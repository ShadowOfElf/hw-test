package main

import (
	"context"
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
	Send(ctx context.Context) error
	Receive(ctx context.Context) error
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

func (c *clientT) Send(ctx context.Context) error {
	end := make(chan struct{})
	go func() {
		writer := &contextAwareWriter{
			ctx:    ctx,
			writer: c.conn,
		}

		_, err := io.Copy(writer, c.in)
		if err != nil && errors.Is(err, context.Canceled) {
			_, _ = fmt.Fprintf(os.Stderr, "Error during copy: %v\n", err)
		}

		close(end)
	}()

	for {
		select {
		case <-ctx.Done():
			_, _ = fmt.Fprintf(os.Stderr, "received stop signal, exiting... \n")
			return nil
		case <-end:
			return nil
		}
	}
}

type contextAwareWriter struct {
	ctx    context.Context
	writer io.Writer
}

func (w *contextAwareWriter) Write(p []byte) (n int, err error) {
	select {
	case <-w.ctx.Done():
		return 0, context.Canceled
	default:
		return w.writer.Write(p)
	}
}

func (c *clientT) Receive(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, err := io.Copy(c.out, c.conn)
			if err != nil {
				return err
			}
		}
	}
}

func (c *clientT) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return errorClose
}
