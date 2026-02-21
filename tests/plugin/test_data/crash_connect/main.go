package main

import (
	"context"
	"log"
	"os"

	"github.com/Southclaws/storyden/sdk/go/storyden"
)

//go:generate ./package.nu

func main() {
	ctx := context.Background()

	log.Println("Crash connect plugin starting...")

	_, err := storyden.New(ctx)
	if err != nil {
		log.Fatalf("Failed to create plugin: %v", err)
	}

	log.Println("Crash connect plugin - CRASHING AFTER WEBSOCKET INIT!")
	os.Exit(42)
}
