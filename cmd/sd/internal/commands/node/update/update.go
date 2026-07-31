package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
)

func New(store *config.Store) cligen.NodeUpdateHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodeUpdateParams) error {
		props, err := buildMutableProps(io, p)
		if err != nil {
			return err
		}

		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		node, err := nodeapi.Update(ctx, client.OpenAPI, p.Slug, props)
		if err != nil {
			return err
		}

		fmt.Fprintf(io.Out, "Updated node: %s (slug: %s)\n", node.Name, node.Slug)

		return nil
	}
}

// changedFields lists which discrete field flags carry a non-default value,
// used only to compose the "cannot combine --json with ..." error message.
// Unlike the original Cobra-backed check (cmd.Flags().Changed), this is a
// value-based proxy: an explicit but empty flag (e.g. --name "") is
// indistinguishable from an unset one. That's an intentional, narrow
// trade-off — nobody sets a discrete field flag to its own zero value on
// purpose while also passing --json.
func changedFields(p cligen.NodeUpdateParams) []string {
	var changed []string
	if p.Name != "" {
		changed = append(changed, "--name")
	}
	if p.NewSlug != "" {
		changed = append(changed, "--new-slug")
	}
	if p.Description != "" {
		changed = append(changed, "--description")
	}
	if p.Content != "" {
		changed = append(changed, "--content")
	}
	if p.ContentFile != "" {
		changed = append(changed, "--content-file")
	}
	if p.Url != "" {
		changed = append(changed, "--url")
	}
	if p.ClearUrl {
		changed = append(changed, "--clear-url")
	}
	if p.HideChildTreeSet {
		changed = append(changed, "--hide-child-tree")
	}
	if len(p.Tags) > 0 {
		changed = append(changed, "--tags")
	}
	return changed
}

func buildMutableProps(io cligen.IO, p cligen.NodeUpdateParams) (openapi.NodeMutableProps, error) {
	if p.Json != "" {
		if changed := changedFields(p); len(changed) > 0 {
			return openapi.NodeMutableProps{}, fmt.Errorf("cannot combine --json with %s", strings.Join(changed, ", "))
		}

		return readJSONProps(p.Json, io.In)
	}

	if p.Url != "" && p.ClearUrl {
		return openapi.NodeMutableProps{}, fmt.Errorf("cannot specify both --url and --clear-url")
	}

	finalContent, err := readContent(p.Content, p.ContentFile, io.In)
	if err != nil {
		return openapi.NodeMutableProps{}, err
	}

	finalContent, err = contentToHTML(finalContent, p.Markdown)
	if err != nil {
		return openapi.NodeMutableProps{}, err
	}

	props := openapi.NodeMutableProps{
		Name:        stringPtr(p.Name),
		Slug:        stringPtr(p.NewSlug),
		Description: stringPtr(p.Description),
		Content:     stringPtr(finalContent),
	}

	if len(p.Tags) > 0 {
		props.Tags = &p.Tags
	}
	if p.Url != "" {
		props.Url.Set(p.Url)
	}
	if p.ClearUrl {
		props.Url.SetNull()
	}
	if p.HideChildTreeSet {
		props.HideChildTree = &p.HideChildTree
	}

	return props, nil
}

func readJSONProps(source string, stdin io.Reader) (openapi.NodeMutableProps, error) {
	var data []byte

	if source == "-" {
		bytes, err := io.ReadAll(stdin)
		if err != nil {
			return openapi.NodeMutableProps{}, fmt.Errorf("failed to read JSON from stdin: %w", err)
		}
		data = bytes
	} else {
		bytes, err := os.ReadFile(source)
		if err != nil {
			return openapi.NodeMutableProps{}, fmt.Errorf("failed to read JSON file: %w", err)
		}
		data = bytes
	}

	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil {
		return openapi.NodeMutableProps{}, fmt.Errorf("invalid node update JSON: %w", err)
	}
	if object == nil {
		return openapi.NodeMutableProps{}, fmt.Errorf("node update JSON must be an object")
	}

	var props openapi.NodeMutableProps
	if err := json.Unmarshal(data, &props); err != nil {
		return openapi.NodeMutableProps{}, fmt.Errorf("invalid node update JSON: %w", err)
	}

	return props, nil
}

// contentToHTML converts Markdown content to HTML when the --markdown flag is
// set, reusing the same converter the backend uses. Without the flag, content
// passes through unchanged because the API already expects HTML.
func contentToHTML(content string, markdown bool) (string, error) {
	if !markdown || content == "" {
		return content, nil
	}

	rt, err := datagraph.NewRichTextFromMarkdown(content)
	if err != nil {
		return "", fmt.Errorf("failed to convert markdown content: %w", err)
	}

	return rt.HTML(), nil
}

func readContent(content string, contentFile string, stdin io.Reader) (string, error) {
	if content != "" && contentFile != "" {
		return "", fmt.Errorf("cannot specify both --content and --content-file")
	}

	if content != "" {
		return content, nil
	}

	if contentFile == "" {
		return "", nil
	}

	if contentFile == "-" {
		bytes, err := io.ReadAll(stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read from stdin: %w", err)
		}

		return string(bytes), nil
	}

	bytes, err := os.ReadFile(contentFile)
	if err != nil {
		return "", fmt.Errorf("failed to read content file: %w", err)
	}

	return string(bytes), nil
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}
