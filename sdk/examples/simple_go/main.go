package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"

	"github.com/Southclaws/storyden/lib/plugin/rpc"
	"github.com/Southclaws/storyden/sdk/go/storyden"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	godotenv.Load(".env")

	rpcURL := os.Getenv("STORYDEN_RPC_URL")

	pl, err := storyden.New(ctx)
	if err != nil {
		log.Fatalf("failed to initialise plugin: %v", err)
	}
	defer func() {
		if err := pl.Shutdown(); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	pl.OnThreadPublished(func(ctx context.Context, event *rpc.EventThreadPublished) error {
		fmt.Printf("thread published: %+v\n", event)
		return nil
	})

	log.Printf("connecting to Storyden RPC: %s", rpcURL)

	if err := pl.Run(ctx); err != nil {
		log.Fatalf("plugin stopped: %v", err)
	}
}
