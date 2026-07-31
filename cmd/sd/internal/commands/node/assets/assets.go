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
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
)

func NewUpload(store *config.Store) cligen.NodeAssetsUploadHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodeAssetsUploadParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		assetName := p.Name
		if assetName == "" {
			assetName = filepath.Base(p.File)
		}

		asset, err := uploadAsset(ctx, client.OpenAPI, p.File, assetName)
		if err != nil {
			return err
		}

		node, err := attachAsset(ctx, client.OpenAPI, p.Slug, asset.Id, p.Primary)
		if err != nil {
			return err
		}

		fmt.Fprintf(io.Out, "Uploaded asset: %s (id: %s)\n", asset.Filename, asset.Id)
		fmt.Fprintf(io.Out, "Attached to node: %s (slug: %s)\n", node.Name, node.Slug)
		return nil
	}
}

func NewPrimarySet(store *config.Store) cligen.NodeAssetsPrimarySetHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodeAssetsPrimarySetParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		node, err := setPrimaryAsset(ctx, client.OpenAPI, p.Slug, p.AssetId)
		if err != nil {
			return err
		}

		fmt.Fprintf(io.Out, "Set primary image for node: %s (slug: %s)\n", node.Name, node.Slug)
		return nil
	}
}

func NewPrimaryClear(store *config.Store) cligen.NodeAssetsPrimaryClearHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodeAssetsPrimaryClearParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		node, err := clearPrimaryAsset(ctx, client.OpenAPI, p.Slug)
		if err != nil {
			return err
		}

		fmt.Fprintf(io.Out, "Cleared primary image for node: %s (slug: %s)\n", node.Name, node.Slug)
		return nil
	}
}

func NewPrimaryDownload(store *config.Store) cligen.NodeAssetsPrimaryDownloadHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodeAssetsPrimaryDownloadParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		node, err := nodeapi.Fetch(ctx, client.OpenAPI, p.Slug)
		if err != nil {
			return err
		}
		if node.PrimaryImage == nil {
			return fmt.Errorf("node has no primary image")
		}

		asset := *node.PrimaryImage
		data, err := downloadAsset(ctx, client.OpenAPI, assetFilename(asset))
		if err != nil {
			return err
		}

		target := p.OutputFile
		if target == "" {
			target = asset.Filename
		}

		if err := writeAssetData(io.Out, data, target, p.Force); err != nil {
			return err
		}

		if target != "-" {
			fmt.Fprintf(io.Out, "Downloaded primary image: %s -> %s\n", asset.Filename, target)
		}
		return nil
	}
}

func NewAdd(store *config.Store) cligen.NodeAssetsAddHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodeAssetsAddParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		node, err := attachAsset(ctx, client.OpenAPI, p.Slug, p.AssetId, p.Primary)
		if err != nil {
			return err
		}

		fmt.Fprintf(io.Out, "Attached asset %s to node: %s (slug: %s)\n", p.AssetId, node.Name, node.Slug)
		return nil
	}
}

func NewRemove(store *config.Store) cligen.NodeAssetsRemoveHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodeAssetsRemoveParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		node, err := removeAsset(ctx, client.OpenAPI, p.Slug, p.AssetId)
		if err != nil {
			return err
		}

		fmt.Fprintf(io.Out, "Removed asset %s from node: %s (slug: %s)\n", p.AssetId, node.Name, node.Slug)
		return nil
	}
}

func NewDownload(store *config.Store) cligen.NodeAssetsDownloadHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodeAssetsDownloadParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		node, err := nodeapi.Fetch(ctx, client.OpenAPI, p.Slug)
		if err != nil {
			return err
		}

		asset, err := findAsset(node.Assets, p.Asset)
		if err != nil {
			return err
		}

		data, err := downloadAsset(ctx, client.OpenAPI, assetFilename(asset))
		if err != nil {
			return err
		}

		target := p.OutputFile
		if target == "" {
			target = asset.Filename
		}

		if err := writeAssetData(io.Out, data, target, p.Force); err != nil {
			return err
		}

		fmt.Fprintf(io.Out, "Downloaded asset: %s -> %s\n", asset.Filename, target)
		return nil
	}
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
