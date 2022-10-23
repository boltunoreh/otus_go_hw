package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connect timeout")

	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("Connection error: %v", err)
	}
	defer client.Close()

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go func(ctx context.Context) {
		defer wg.Done()
	OUTER:
		for {
			select {
			case <-ctx.Done():
				break OUTER
			default:
				err := client.Receive()
				if err != nil {
					log.Fatalf("Receive error: %v", err)
				}
			}
		}

		fmt.Fprintf(os.Stderr, "...Connection was closed by peer\n")
		cancel()
	}(ctx)

	go func(ctx context.Context) {
		defer wg.Done()
	OUTER:
		for {
			select {
			case <-ctx.Done():
				break OUTER
			default:
				err := client.Send()
				if err != nil {
					log.Fatalf("Send error: %v", err)
				}
			}
		}

		fmt.Fprintf(os.Stderr, "...EOF\n")
		cancel()
	}(ctx)

	wg.Wait()
	<-ctx.Done()
}
