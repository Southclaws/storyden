package pluginnew

import (
	"context"
	"testing"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
)

func TestNewRejectsDirectoryWithoutSlug(t *testing.T) {
	handler := New()

	err := handler(context.Background(), &cobra.Command{}, cligen.IO{}, cligen.PluginDevNewParams{Directory: "."})
	if err == nil {
		t.Fatal("expected error")
	}
}
