package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Southclaws/storyden/lib/plugin/rpc"
	"github.com/Southclaws/storyden/sdk/go/storyden"
	"github.com/rs/xid"
)

//go:generate ./package.nu

func main() {
	outputDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	log.Printf("Output directory: %s", outputDir)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p, err := storyden.New(ctx)
	if err != nil {
		log.Fatalf("Failed to create plugin: %v", err)
	}

	p.OnThreadPublished(func(ctx context.Context, event *rpc.EventThreadPublished) error {
		return handleEvent(outputDir, "thread_published", xid.ID(event.ID), event)
	})

	p.OnThreadUnpublished(func(ctx context.Context, event *rpc.EventThreadUnpublished) error {
		return handleEvent(outputDir, "thread_unpublished", xid.ID(event.ID), event)
	})

	p.OnThreadUpdated(func(ctx context.Context, event *rpc.EventThreadUpdated) error {
		return handleEvent(outputDir, "thread_updated", xid.ID(event.ID), event)
	})

	p.OnThreadDeleted(func(ctx context.Context, event *rpc.EventThreadDeleted) error {
		return handleEvent(outputDir, "thread_deleted", xid.ID(event.ID), event)
	})

	p.OnThreadReplyCreated(func(ctx context.Context, event *rpc.EventThreadReplyCreated) error {
		return handleEvent(outputDir, "thread_reply_created", xid.ID(event.ReplyID), event)
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	log.Println("Plugin started, waiting for events...")

	if err := p.Run(ctx); err != nil {
		log.Printf("Plugin error: %v", err)
	}

	log.Println("Plugin stopped")
}

func handleEvent(outputDir, eventType string, id xid.ID, event any) error {
	log.Printf("Received %s event: %s", eventType, id)

	filename := filepath.Join(outputDir, fmt.Sprintf("%s.json", id))
	data, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err := os.WriteFile(filename, data, 0o644); err != nil {
		return fmt.Errorf("failed to write event file: %w", err)
	}

	log.Printf("Wrote %s event to %s", eventType, filename)
	return nil
}
