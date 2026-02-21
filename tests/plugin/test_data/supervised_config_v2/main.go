package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	sdk "github.com/Southclaws/storyden/sdk/go/storyden"
)

//go:generate ./package.nu

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	plugin, err := sdk.New(ctx)
	if err != nil {
		log.Fatalf("failed to create plugin: %v", err)
	}

	plugin.OnConfigure(func(_ context.Context, config map[string]any) error {
		return writeConfiguredFile(config)
	})

	if err := plugin.Run(ctx); err != nil {
		log.Printf("plugin stopped with error: %v", err)
	}
}

func writeConfiguredFile(config map[string]any) error {
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile("configured.json", b, 0o644); err != nil {
		return fmt.Errorf("failed to write configured.json: %w", err)
	}

	return nil
}
