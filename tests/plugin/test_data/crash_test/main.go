package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Southclaws/storyden/lib/plugin/rpc"
	"github.com/Southclaws/storyden/sdk/go/storyden"
)

//go:generate ./package.nu

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p, err := storyden.New(ctx)
	if err != nil {
		log.Fatalf("Failed to create plugin: %v", err)
	}

	p.OnThreadPublished(func(ctx context.Context, event *rpc.EventThreadPublished) error {
		log.Printf("Received thread published event: %s - CRASHING NOW!", event.ID)
		os.Exit(42)
		return nil
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	log.Println("Crash test plugin started, waiting for events...")

	if err := p.Run(ctx); err != nil {
		log.Printf("Plugin error: %v", err)
	}

	log.Println("Plugin stopped")
}
