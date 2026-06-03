package pluginapi

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"

	pluginresource "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

const (
	ManifestFilename = "manifest.yaml"
	RPCURLEnvName    = "STORYDEN_RPC_URL"
)

type ManifestFile struct {
	Path     string
	Manifest rpc.Manifest
}

type ExternalPlugin struct {
	ID    string
	Token string
}

type PackageArchive struct {
	Manifest rpc.Manifest
	Bytes    []byte
	Files    []string
}

type excludePath struct {
	path string
	dir  bool
}

func ReadManifest(path string) (*ManifestFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	manifest, err := rpc.ManifestFromMap(raw)
	if err != nil {
		return nil, fmt.Errorf("validate manifest: %w", err)
	}

	return &ManifestFile{Path: path, Manifest: *manifest}, nil
}

func ReadProjectManifest(dir string, manifestPath string) (*ManifestFile, error) {
	path := manifestPath
	if path == "" {
		path = ManifestFilename
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(dir, path)
	}
	return ReadManifest(path)
}

func WriteNewManifest(out io.Writer, dir string, manifest rpc.Manifest, force bool) error {
	if err := manifest.Validate(); err != nil {
		return fmt.Errorf("validate manifest: %w", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	path := filepath.Join(dir, ManifestFilename)
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists; use --force to overwrite", path)
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}

	fmt.Fprintf(out, "Created plugin manifest at %s\n", path)
	return nil
}

func BuildPackage(ctx context.Context, dir string, manifestPath string, excludePaths ...string) (*PackageArchive, error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	mf, err := ReadProjectManifest(root, manifestPath)
	if err != nil {
		return nil, err
	}

	manifestJSON, err := json.MarshalIndent(mf.Manifest, "", "  ")
	if err != nil {
		return nil, err
	}
	manifestJSON = append(manifestJSON, '\n')

	excludes := []excludePath{}
	for _, path := range excludePaths {
		if path == "" {
			continue
		}
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, path)
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		abs = filepath.Clean(abs)
		info, err := os.Stat(abs)
		excludes = append(excludes, excludePath{
			path: abs,
			dir:  err == nil && info.IsDir(),
		})
	}

	type archiveFile struct {
		abs string
		rel string
	}
	files := []archiveFile{}
	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}

		abs := filepath.Clean(path)
		if excluded := excludedArchivePath(abs, excludes); excluded {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}

		name := d.Name()
		if d.IsDir() {
			switch name {
			case ".git", "node_modules", ".next":
				return filepath.SkipDir
			}
			return nil
		}
		if name == ".DS_Store" || name == pluginresource.ArchiveManifestFileName {
			return nil
		}

		if sameFilePath(abs, mf.Path) {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, archiveFile{abs: abs, rel: filepath.ToSlash(rel)})
		return nil
	}); err != nil {
		return nil, err
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].rel < files[j].rel
	})

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	manifestHeader := &zip.FileHeader{
		Name:   pluginresource.ArchiveManifestFileName,
		Method: zip.Deflate,
	}
	manifestHeader.SetMode(0o644)
	manifestWriter, err := zw.CreateHeader(manifestHeader)
	if err != nil {
		return nil, err
	}
	if _, err := manifestWriter.Write(manifestJSON); err != nil {
		return nil, err
	}

	written := []string{pluginresource.ArchiveManifestFileName}
	for _, file := range files {
		linkInfo, err := os.Lstat(file.abs)
		if err != nil {
			return nil, err
		}
		if linkInfo.Mode()&os.ModeSymlink != 0 {
			continue
		}
		info, err := os.Stat(file.abs)
		if err != nil {
			return nil, err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return nil, err
		}
		header.Name = file.rel
		header.Method = zip.Deflate

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return nil, err
		}
		reader, err := os.Open(file.abs)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(writer, reader); err != nil {
			_ = reader.Close()
			return nil, err
		}
		if err := reader.Close(); err != nil {
			return nil, err
		}
		written = append(written, file.rel)
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	archive := pluginresource.Binary(buf.Bytes())
	validated, err := archive.Validate(ctx)
	if err != nil {
		return nil, fmt.Errorf("validate package: %w", err)
	}
	if err := validated.Metadata.Validate(); err != nil {
		return nil, fmt.Errorf("validate package manifest: %w", err)
	}

	return &PackageArchive{
		Manifest: mf.Manifest,
		Bytes:    buf.Bytes(),
		Files:    written,
	}, nil
}

func excludedArchivePath(path string, excludes []excludePath) bool {
	for _, exclude := range excludes {
		if path == exclude.path {
			return true
		}
		if exclude.dir {
			prefix := exclude.path + string(filepath.Separator)
			if strings.HasPrefix(path, prefix) {
				return true
			}
		}
	}
	return false
}

func WritePackageFile(path string, pkg *PackageArchive, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists; use --force to overwrite", path)
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, pkg.Bytes, 0o644)
}

func DefaultPackagePath(dir string, manifest rpc.Manifest) string {
	return filepath.Join(dir, Slugify(manifest.ID)+".zip")
}

func sameFilePath(a string, b string) bool {
	aa, err := filepath.Abs(a)
	if err != nil {
		return filepath.Clean(a) == filepath.Clean(b)
	}
	bb, err := filepath.Abs(b)
	if err != nil {
		return filepath.Clean(a) == filepath.Clean(b)
	}
	return filepath.Clean(aa) == filepath.Clean(bb)
}

func EnsureExternalPlugin(ctx context.Context, client *openapi.ClientWithResponses, manifest rpc.Manifest, requestedID string, noUpdate bool) (*ExternalPlugin, error) {
	if requestedID != "" {
		var plugin *openapi.Plugin
		var err error
		if !noUpdate {
			plugin, err = UpdateManifest(ctx, client, requestedID, manifest)
			if err != nil {
				return nil, err
			}
		} else {
			plugin, err = GetPlugin(ctx, client, requestedID)
			if err != nil {
				return nil, err
			}
		}
		return ExternalPluginFromAPI(*plugin)
	}

	plugins, err := ListPlugins(ctx, client)
	if err != nil {
		return nil, err
	}
	for _, pl := range plugins {
		if id, _ := pl.Manifest["id"].(string); id == manifest.ID {
			if PluginMode(pl) != string(openapi.External) {
				return nil, fmt.Errorf("plugin manifest id %q is already installed as a supervised plugin; pass --instance-id for an external installation", manifest.ID)
			}
			if !noUpdate {
				updated, err := UpdateManifest(ctx, client, string(pl.Id), manifest)
				if err != nil {
					return nil, err
				}
				return ExternalPluginFromAPI(*updated)
			}
			return ExternalPluginFromAPI(pl)
		}
	}

	body := openapi.PluginInitialProps{}
	if err := body.FromPluginInitialExternal(openapi.PluginInitialExternal{
		Mode:     openapi.External,
		Manifest: openapi.PluginManifest(manifest.ToMap()),
	}); err != nil {
		return nil, err
	}

	response, err := client.PluginAddWithResponse(ctx, body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, output.RequestErrorWithMessages("plugin add request", response, response.Body, output.UnauthorizedMessage("plugin add request"))
	}
	return ExternalPluginFromAPI(*response.JSON200)
}

func InstallSupervisedPackage(ctx context.Context, client *openapi.ClientWithResponses, pkg *PackageArchive) (*openapi.Plugin, bool, error) {
	plugins, err := ListPlugins(ctx, client)
	if err != nil {
		return nil, false, err
	}

	for _, pl := range plugins {
		if id, _ := pl.Manifest["id"].(string); id == pkg.Manifest.ID {
			if PluginMode(pl) != string(openapi.Supervised) {
				return nil, false, fmt.Errorf("plugin manifest id %q is already installed as an external plugin", pkg.Manifest.ID)
			}

			response, err := client.PluginUpdatePackageWithBodyWithResponse(ctx, openapi.PluginIDParam(pl.Id), "application/zip", bytes.NewReader(pkg.Bytes))
			if err != nil {
				return nil, false, err
			}
			if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
				return nil, false, output.RequestErrorWithMessages("plugin package update request", response, response.Body, output.UnauthorizedMessage("plugin package update request"))
			}
			return (*openapi.Plugin)(response.JSON200), true, nil
		}
	}

	response, err := client.PluginAddWithBodyWithResponse(ctx, "application/zip", bytes.NewReader(pkg.Bytes))
	if err != nil {
		return nil, false, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, false, output.RequestErrorWithMessages("plugin package install request", response, response.Body, output.UnauthorizedMessage("plugin package install request"))
	}
	return (*openapi.Plugin)(response.JSON200), false, nil
}

func ListPlugins(ctx context.Context, client *openapi.ClientWithResponses) ([]openapi.Plugin, error) {
	response, err := client.PluginListWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, output.RequestErrorWithMessages("plugin list request", response, response.Body, output.UnauthorizedMessage("plugin list request"))
	}
	return response.JSON200.Plugins, nil
}

func ExternalPluginFromAPI(p openapi.Plugin) (*ExternalPlugin, error) {
	if PluginMode(p) != string(openapi.External) {
		return nil, fmt.Errorf("plugin %s is not an external plugin", p.Id)
	}
	props, err := p.Connection.AsPluginExternalProps()
	if err != nil {
		return nil, fmt.Errorf("plugin %s is not an external plugin: %w", p.Id, err)
	}
	if props.Token == "" {
		return nil, fmt.Errorf("plugin %s did not include an external RPC token", p.Id)
	}
	return &ExternalPlugin{ID: string(p.Id), Token: props.Token}, nil
}

func GetPlugin(ctx context.Context, client *openapi.ClientWithResponses, id string) (*openapi.Plugin, error) {
	response, err := client.PluginGetWithResponse(ctx, openapi.PluginIDParam(id))
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, output.RequestErrorWithMessages("plugin get request", response, response.Body, output.UnauthorizedMessage("plugin get request"))
	}
	return (*openapi.Plugin)(response.JSON200), nil
}

func DeletePlugin(ctx context.Context, client *openapi.ClientWithResponses, id string) error {
	response, err := client.PluginDeleteWithResponse(ctx, openapi.PluginIDParam(id))
	if err != nil {
		return err
	}
	if response.StatusCode() != http.StatusNoContent {
		return output.RequestErrorWithMessages("plugin delete request", response, response.Body, output.UnauthorizedMessage("plugin delete request"))
	}
	return nil
}

func SetActiveState(ctx context.Context, client *openapi.ClientWithResponses, id string, state openapi.PluginActiveState) (*openapi.Plugin, error) {
	response, err := client.PluginSetActiveStateWithResponse(ctx, openapi.PluginIDParam(id), openapi.PluginSetActiveStateJSONRequestBody{Active: state})
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, output.RequestErrorWithMessages("plugin active-state request", response, response.Body, output.UnauthorizedMessage("plugin active-state request"))
	}
	return (*openapi.Plugin)(response.JSON200), nil
}

func UpdateManifest(ctx context.Context, client *openapi.ClientWithResponses, id string, manifest rpc.Manifest) (*openapi.Plugin, error) {
	response, err := client.PluginUpdateManifestWithResponse(ctx, openapi.PluginIDParam(id), openapi.PluginUpdateManifestJSONRequestBody(manifest.ToMap()))
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, output.RequestErrorWithMessages("plugin manifest update request", response, response.Body, output.UnauthorizedMessage("plugin manifest update request"))
	}
	return (*openapi.Plugin)(response.JSON200), nil
}

func CycleToken(ctx context.Context, client *openapi.ClientWithResponses, id string) (string, error) {
	response, err := client.PluginCycleTokenWithResponse(ctx, openapi.PluginIDParam(id))
	if err != nil {
		return "", err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return "", output.RequestErrorWithMessages("plugin token request", response, response.Body, output.UnauthorizedMessage("plugin token request"))
	}
	return response.JSON200.Token, nil
}

func ExternalRPCURL(endpoint string, token string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	case "ws", "wss":
	default:
		return "", fmt.Errorf("unsupported Storyden endpoint scheme %q", u.Scheme)
	}
	u.Path = "/rpc"
	u.RawPath = ""
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func CommandFromManifest(manifest rpc.Manifest, override []string) (string, []string, error) {
	if len(override) > 0 {
		return override[0], override[1:], nil
	}
	if strings.TrimSpace(manifest.Command) == "" {
		return "", nil, fmt.Errorf("manifest command is empty; pass a command after --")
	}
	return manifest.Command, manifest.Args, nil
}

func PluginMode(p openapi.Plugin) string {
	mode, err := p.Connection.Discriminator()
	if err != nil {
		return ""
	}
	return mode
}

func PluginStatus(p openapi.Plugin) string {
	state, err := p.Status.Discriminator()
	if err != nil {
		return ""
	}
	return state
}

func DefaultAuthor() string {
	if user := os.Getenv("USER"); user != "" {
		return user
	}
	if user := os.Getenv("USERNAME"); user != "" {
		return user
	}
	return "you"
}

func Slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		valid := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if valid {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash && b.Len() > 0 {
			b.WriteRune('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func Titleize(value string) string {
	value = strings.ReplaceAll(value, "-", " ")
	value = strings.ReplaceAll(value, "_", " ")
	parts := strings.Fields(value)
	for i, part := range parts {
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	if len(parts) == 0 {
		return "My Plugin"
	}
	return strings.Join(parts, " ")
}
