package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Operation timeout")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		_, _ = fmt.Fprint(os.Stderr, "Usage: go-telnet --timeout=TIMEOUT host port")
		os.Exit(1)
	}
	host := args[0]
	port := args[1]

	if host == "" || port == "" {
		_, _ = fmt.Fprint(os.Stderr, "Usage: go-telnet wrong host or port")
		os.Exit(1)
	}
	address := net.JoinHostPort(host, port)

	clientTelnet := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	TelnetWork(clientTelnet)
	_, _ = fmt.Fprint(os.Stderr, "telnet work ended")
}

func TelnetWork(clientTelnet TelnetClient) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err := clientTelnet.Connect()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "Connect to host error")
	}
	defer func() {
		err := clientTelnet.Close()
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, "close connected to host error")
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := clientTelnet.Send(ctx); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "send msg error: %s\n", err)
		}
	}()

	go func() {
		if err := clientTelnet.Receive(ctx); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "read msg from server error: %s\n", err)
		}
	}()

	wg.Wait()
}
