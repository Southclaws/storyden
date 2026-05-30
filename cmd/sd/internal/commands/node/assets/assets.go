package assets

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
)

type AssetsCommand *cobra.Command

func New(store *config.Store) AssetsCommand {
	command := &cobra.Command{
		Use:   "assets",
		Short: "Upload, download, attach, and remove node assets",
		Long: `# Node Assets

Work with files attached to nodes.

Assets can be uploaded from local files, attached to a node, removed from a node, or downloaded from an existing node attachment.

Primary images are handled separately with ` + "`sd node assets primary`" + `. They are used as cover or hero images on pages and may not appear in the normal attached assets list.
`,
	}

	command.AddCommand(newUploadCommand(store))
	command.AddCommand(newAddCommand(store))
	command.AddCommand(newRemoveCommand(store))
	command.AddCommand(newDownloadCommand(store))
	command.AddCommand(newPrimaryCommand(store))

	help.SetupMarkdownHelp(command)

	return AssetsCommand(command)
}

func newUploadCommand(store *config.Store) *cobra.Command {
	var name string
	var primary bool

	command := &cobra.Command{
		Use:   "upload <slug> <file>",
		Short: "Upload a file and attach it to a node",
		Long: `# Upload a Node Asset

Upload a local file to Storyden and attach it to a node.

## Examples

Upload an image:
~~~bash
sd node assets upload docs ./diagram.png
~~~

Upload and mark as the primary image:
~~~bash
sd node assets upload docs ./cover.png --primary
~~~
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			assetName := name
			if assetName == "" {
				assetName = filepath.Base(args[1])
			}

			asset, err := uploadAsset(cmd.Context(), client.OpenAPI, args[1], assetName)
			if err != nil {
				return err
			}

			node, err := attachAsset(cmd.Context(), client.OpenAPI, args[0], asset.Id, primary)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Uploaded asset: %s (id: %s)\n", asset.Filename, asset.Id)
			fmt.Fprintf(cmd.OutOrStdout(), "Attached to node: %s (slug: %s)\n", node.Name, node.Slug)
			return nil
		},
	}

	command.Flags().StringVar(&name, "name", "", "Asset filename to store")
	command.Flags().BoolVar(&primary, "primary", false, "Set the asset as the node's primary image")
	help.SetupMarkdownHelp(command)

	return command
}

func newPrimaryCommand(store *config.Store) *cobra.Command {
	command := &cobra.Command{
		Use:   "primary",
		Short: "Set, clear, or download the node primary image",
		Long: `# Node Primary Image

Work with a node's primary image. This is the cover or hero image used on pages and is separate from the normal attached assets list.
`,
	}

	command.AddCommand(newPrimarySetCommand(store))
	command.AddCommand(newPrimaryClearCommand(store))
	command.AddCommand(newPrimaryDownloadCommand(store))

	help.SetupMarkdownHelp(command)

	return command
}

func newPrimarySetCommand(store *config.Store) *cobra.Command {
	command := &cobra.Command{
		Use:   "set <slug> <asset-id>",
		Short: "Set the node primary image",
		Long: `# Set Node Primary Image

Set the asset used as a node's cover or hero image.

## Examples

Set by asset ID:
~~~bash
sd node assets primary set docs d8cg...
~~~
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			node, err := setPrimaryAsset(cmd.Context(), client.OpenAPI, args[0], args[1])
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Set primary image for node: %s (slug: %s)\n", node.Name, node.Slug)
			return nil
		},
	}

	help.SetupMarkdownHelp(command)

	return command
}

func newPrimaryClearCommand(store *config.Store) *cobra.Command {
	command := &cobra.Command{
		Use:   "clear <slug>",
		Short: "Clear the node primary image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			node, err := clearPrimaryAsset(cmd.Context(), client.OpenAPI, args[0])
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Cleared primary image for node: %s (slug: %s)\n", node.Name, node.Slug)
			return nil
		},
	}

	help.SetupMarkdownHelp(command)

	return command
}

func newPrimaryDownloadCommand(store *config.Store) *cobra.Command {
	var output string
	var force bool

	command := &cobra.Command{
		Use:   "download <slug>",
		Short: "Download the node primary image",
		Long: `# Download Node Primary Image

Download the asset used as a node's cover or hero image.

## Examples

Download to the asset filename:
~~~bash
sd node assets primary download docs
~~~

Download to a specific file:
~~~bash
sd node assets primary download docs --output cover.png
~~~
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			node, err := nodeapi.Fetch(cmd.Context(), client.OpenAPI, args[0])
			if err != nil {
				return err
			}
			if node.PrimaryImage == nil {
				return fmt.Errorf("node has no primary image")
			}

			asset := *node.PrimaryImage
			data, err := downloadAsset(cmd.Context(), client.OpenAPI, assetFilename(asset))
			if err != nil {
				return err
			}

			target := output
			if target == "" {
				target = asset.Filename
			}

			if err := writeAssetData(cmd.OutOrStdout(), data, target, force); err != nil {
				return err
			}

			if target != "-" {
				fmt.Fprintf(cmd.OutOrStdout(), "Downloaded primary image: %s -> %s\n", asset.Filename, target)
			}
			return nil
		},
	}

	command.Flags().StringVarP(&output, "output", "o", "", "Output file path (use - for stdout)")
	command.Flags().BoolVar(&force, "force", false, "Overwrite the output file if it exists")
	help.SetupMarkdownHelp(command)

	return command
}

func newAddCommand(store *config.Store) *cobra.Command {
	var primary bool

	command := &cobra.Command{
		Use:   "add <slug> <asset-id>",
		Short: "Attach an existing asset to a node",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			node, err := attachAsset(cmd.Context(), client.OpenAPI, args[0], args[1], primary)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Attached asset %s to node: %s (slug: %s)\n", args[1], node.Name, node.Slug)
			return nil
		},
	}

	command.Flags().BoolVar(&primary, "primary", false, "Set the asset as the node's primary image")
	help.SetupMarkdownHelp(command)

	return command
}

func newRemoveCommand(store *config.Store) *cobra.Command {
	command := &cobra.Command{
		Use:   "remove <slug> <asset-id>",
		Short: "Remove an asset from a node",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			node, err := removeAsset(cmd.Context(), client.OpenAPI, args[0], args[1])
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Removed asset %s from node: %s (slug: %s)\n", args[1], node.Name, node.Slug)
			return nil
		},
	}

	help.SetupMarkdownHelp(command)

	return command
}

func newDownloadCommand(store *config.Store) *cobra.Command {
	var output string
	var force bool

	command := &cobra.Command{
		Use:   "download <slug> <asset>",
		Short: "Download an attached node asset",
		Long: `# Download a Node Asset

Download an asset that is already attached to a node. The asset can be identified by asset ID, filename, or API path.

## Examples

Download by asset ID:
~~~bash
sd node assets download docs d8cg... --output cover.png
~~~

Download by filename:
~~~bash
sd node assets download docs cover.png
~~~
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			node, err := nodeapi.Fetch(cmd.Context(), client.OpenAPI, args[0])
			if err != nil {
				return err
			}

			asset, err := findAsset(node.Assets, args[1])
			if err != nil {
				return err
			}

			data, err := downloadAsset(cmd.Context(), client.OpenAPI, assetFilename(asset))
			if err != nil {
				return err
			}

			target := output
			if target == "" {
				target = asset.Filename
			}

			if err := writeAssetData(cmd.OutOrStdout(), data, target, force); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Downloaded asset: %s -> %s\n", asset.Filename, target)
			return nil
		},
	}

	command.Flags().StringVarP(&output, "output", "o", "", "Output file path (use - for stdout)")
	command.Flags().BoolVar(&force, "force", false, "Overwrite the output file if it exists")
	help.SetupMarkdownHelp(command)

	return command
}

func uploadAsset(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	filePath string,
	name string,
) (*openapi.Asset, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open asset file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to inspect asset file: %w", err)
	}

	contentType := detectContentType(file, name)
	params := &openapi.AssetUploadParams{
		Filename:      &name,
		ContentLength: stat.Size(),
	}

	response, err := client.AssetUploadWithBodyWithResponse(ctx, params, contentType, file, func(ctx context.Context, req *http.Request) error {
		req.ContentLength = stat.Size()
		return nil
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, assetUploadError(response)
	}

	return response.JSON200, nil
}

func setPrimaryAsset(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	assetID string,
) (*openapi.NodeWithChildren, error) {
	props := openapi.NodeMutableProps{}
	props.PrimaryImageAssetId.Set(assetID)

	return nodeapi.Update(ctx, client, slug, props)
}

func clearPrimaryAsset(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
) (*openapi.NodeWithChildren, error) {
	props := openapi.NodeMutableProps{}
	props.PrimaryImageAssetId.SetNull()

	return nodeapi.Update(ctx, client, slug, props)
}

func attachAsset(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	assetID string,
	primary bool,
) (*openapi.NodeWithChildren, error) {
	response, err := client.NodeAddAssetWithResponse(ctx, slug, assetID)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, nodeAddAssetError(response)
	}

	node := response.JSON200
	if primary {
		props := openapi.NodeMutableProps{}
		props.PrimaryImageAssetId.Set(assetID)

		updated, err := nodeapi.Update(ctx, client, slug, props)
		if err != nil {
			return nil, err
		}

		node = updated
	}

	return node, nil
}

func removeAsset(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	assetID string,
) (*openapi.NodeWithChildren, error) {
	response, err := client.NodeRemoveAssetWithResponse(ctx, slug, assetID)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, nodeRemoveAssetError(response)
	}

	return response.JSON200, nil
}

func downloadAsset(ctx context.Context, client *openapi.ClientWithResponses, filename string) ([]byte, error) {
	response, err := client.AssetGetWithResponse(ctx, filename)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		return nil, assetGetError(response)
	}

	return response.Body, nil
}

func writeAssetData(out io.Writer, data []byte, target string, force bool) error {
	if target == "-" {
		_, err := out.Write(data)
		return err
	}

	if !force {
		if _, err := os.Stat(target); err == nil {
			return fmt.Errorf("output file already exists: %s (use --force to overwrite)", target)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("failed to inspect output file: %w", err)
		}
	}

	if err := os.WriteFile(target, data, 0o644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func detectContentType(file *os.File, name string) string {
	if contentType := mime.TypeByExtension(filepath.Ext(name)); contentType != "" {
		return contentType
	}

	var header [512]byte
	n, _ := file.Read(header[:])
	_, _ = file.Seek(0, io.SeekStart)

	if n > 0 {
		return http.DetectContentType(header[:n])
	}

	return "application/octet-stream"
}

func findAsset(assets openapi.AssetList, selector string) (openapi.Asset, error) {
	for _, asset := range assets {
		if assetMatches(asset, selector) {
			return asset, nil
		}
	}

	return openapi.Asset{}, fmt.Errorf("asset not attached to node: %s", selector)
}

func assetMatches(asset openapi.Asset, selector string) bool {
	filename := assetFilename(asset)

	return selector == asset.Id ||
		selector == asset.Filename ||
		selector == asset.Path ||
		selector == filename
}

func assetFilename(asset openapi.Asset) string {
	if asset.Path == "" {
		return asset.Filename
	}

	return path.Base(strings.TrimSuffix(asset.Path, "/"))
}

func assetUploadError(response *openapi.AssetUploadResponse) error {
	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("asset upload request failed: %s: %s", response.Status(), body)
	}
	return fmt.Errorf("asset upload request failed: %s", response.Status())
}

func assetGetError(response *openapi.AssetGetResponse) error {
	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("asset download request failed: %s: %s", response.Status(), body)
	}
	return fmt.Errorf("asset download request failed: %s", response.Status())
}

func nodeAddAssetError(response *openapi.NodeAddAssetResponse) error {
	if response.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("node or asset not found")
	}
	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("node asset request was not authorised; run sd auth login again")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("node asset request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("node asset request failed: %s", response.Status())
}

func nodeRemoveAssetError(response *openapi.NodeRemoveAssetResponse) error {
	if response.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("node or asset not found")
	}
	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("node asset request was not authorised; run sd auth login again")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("node asset request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("node asset request failed: %s", response.Status())
}
