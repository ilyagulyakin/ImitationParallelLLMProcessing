package main

import (
	"ImitationParallelLLMProcessing/agent"
	server "ImitationParallelLLMProcessing/sever"
	"context"
	"log"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch := make(chan agent.Ticket)
	go Run(ctx, ch)

	srv := server.NewServer(ch)
	serverErr := make(chan error, 1)
	go func() { serverErr <- srv.Run(":9091") }()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		return nil
	}
}
