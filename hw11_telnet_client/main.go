package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Operation timeout")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("Usage: go-telnet --timeout=TIMEOUT host port")
	}
	host := args[0]
	port := args[1]

	if host == "" || port == "" {
		log.Fatal("Usage: go-telnet wrong host or port")
	}
	address := net.JoinHostPort(host, port)

	clientTelnet := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := TelnetWork(clientTelnet); err != nil {
		log.Fatal("Usage: go-telnet error in work", err)
	}
}
