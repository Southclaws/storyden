package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/oapi-codegen/nullable"
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	commandcontent "github.com/Southclaws/storyden/cmd/sd/internal/content"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
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

		node, err := nodeapi.Update(ctx, client.OpenAPI, p.NodeSlug, props)
		if err != nil {
			return err
		}

		if p.Output == "json" {
			return output.JSON(io.Out, node)
		}

		fmt.Fprintf(io.Out, "Updated page: %s (slug: %s)\n", node.Name, node.Slug)

		return nil
	}
}

func changedFields(p cligen.NodeUpdateParams) []string {
	var changed []string
	if p.NameSet {
		changed = append(changed, "--name")
	}
	if p.SlugSet {
		changed = append(changed, "--slug")
	}
	if p.DescriptionSet {
		changed = append(changed, "--description")
	}
	if p.ContentSet {
		changed = append(changed, "--content")
	}
	if p.ContentFileSet {
		changed = append(changed, "--content-file")
	}
	if p.UrlSet {
		changed = append(changed, "--url")
	}
	if p.ClearUrlSet {
		changed = append(changed, "--clear-url")
	}
	if p.HideChildTreeSet {
		changed = append(changed, "--hide-child-tree")
	}
	if p.TagsSet {
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

	if p.UrlSet && p.ClearUrl {
		return openapi.NodeMutableProps{}, fmt.Errorf("cannot specify both --url and --clear-url")
	}

	finalContent, err := commandcontent.Read(p.Content, p.ContentFile, io.In)
	if err != nil {
		return openapi.NodeMutableProps{}, err
	}

	finalContent, err = commandcontent.ToHTML(finalContent, p.Markdown)
	if err != nil {
		return openapi.NodeMutableProps{}, err
	}

	props := openapi.NodeMutableProps{
		Name:        stringPtr(p.Name),
		Slug:        stringPtr(p.Slug),
		Description: stringPtr(p.Description),
		Content:     stringPtr(finalContent),
	}

	if p.TagsSet {
		props.Tags = &p.Tags
	}
	if p.UrlSet {
		props.Url = nullable.NewNullableWithValue(p.Url)
	}
	if p.ClearUrl {
		props.Url = nullable.NewNullNullable[string]()
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
		return openapi.NodeMutableProps{}, fmt.Errorf("invalid page update JSON: %w", err)
	}
	if object == nil {
		return openapi.NodeMutableProps{}, fmt.Errorf("page update JSON must be an object")
	}

	var props openapi.NodeMutableProps
	if err := json.Unmarshal(data, &props); err != nil {
		return openapi.NodeMutableProps{}, fmt.Errorf("invalid page update JSON: %w", err)
	}

	return props, nil
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}
