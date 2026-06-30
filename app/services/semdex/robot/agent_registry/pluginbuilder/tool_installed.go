package pluginbuilder

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/rs/xid"
	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool/functiontool"

	pluginresource "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
)

type InstalledPluginSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type InstalledPluginsResult struct {
	Plugins []InstalledPluginSummary `json:"plugins"`
}

type ImportInstalledPluginInput struct {
	ID string `json:"id" jsonschema:"Installed supervised plugin ID selected from plugin_installed_list"`
}

type ImportInstalledPluginResult struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Files       int    `json:"files"`
}

type importedPluginArchive struct {
	Manifest *pluginresource.Validated
	Files    int
}

type archiveImportFile struct {
	Path string
	Data []byte
}

func (a *Agent) addInstalledPluginTools(add toolAdder) error {
	if err := add(functiontool.New(functiontool.Config{
		Name:        "plugin_installed_list",
		Description: "List installed supervised plugins that can be imported into the active workspace for editing.",
	}, func(ctx adkagent.Context, args struct{}) (InstalledPluginsResult, error) {
		return a.ListInstalledPlugins(ctx)
	})); err != nil {
		return err
	}

	return add(functiontool.New(functiontool.Config{
		Name:        "plugin_workspace_import_installation",
		Description: "Import one installed supervised plugin archive into an empty active workspace for this chat to edit.",
	}, func(ctx adkagent.Context, args ImportInstalledPluginInput) (ImportInstalledPluginResult, error) {
		return a.ImportInstalledPlugin(ctx, args)
	}))
}

func (a *Agent) ListInstalledPlugins(ctx context.Context) (InstalledPluginsResult, error) {
	if a.installer == nil {
		return InstalledPluginsResult{}, fmt.Errorf("plugin installer is not configured")
	}

	records, err := a.installer.List(ctx)
	if err != nil {
		return InstalledPluginsResult{}, err
	}

	plugins := []InstalledPluginSummary{}
	for _, rec := range records {
		if rec == nil || !rec.Mode.Supervised() {
			continue
		}
		plugins = append(plugins, InstalledPluginSummary{
			ID:          rec.InstallationID.String(),
			Name:        rec.Manifest.Name,
			Description: rec.Manifest.Description,
		})
	}
	sort.Slice(plugins, func(i, j int) bool {
		left := strings.ToLower(plugins[i].Name)
		right := strings.ToLower(plugins[j].Name)
		if left == right {
			return plugins[i].ID < plugins[j].ID
		}
		return left < right
	})

	return InstalledPluginsResult{Plugins: plugins}, nil
}

func (a *Agent) ImportInstalledPlugin(ctx context.Context, in ImportInstalledPluginInput) (ImportInstalledPluginResult, error) {
	if a.installer == nil {
		return ImportInstalledPluginResult{}, fmt.Errorf("plugin installer is not configured")
	}
	if a.reader == nil {
		return ImportInstalledPluginResult{}, fmt.Errorf("plugin reader is not configured")
	}

	rawID := strings.TrimSpace(in.ID)
	parsed, err := xid.FromString(rawID)
	if err != nil {
		return ImportInstalledPluginResult{}, fmt.Errorf("invalid installed plugin id: %w", err)
	}
	id := pluginresource.InstallationID(parsed)

	record, err := a.installer.Get(ctx, id)
	if err != nil {
		return ImportInstalledPluginResult{}, err
	}
	if record == nil || !record.Mode.Supervised() {
		return ImportInstalledPluginResult{}, fmt.Errorf("only supervised plugins can be imported for editing")
	}

	workspace, err := a.Workspace(ctx)
	if err != nil {
		return ImportInstalledPluginResult{}, err
	}
	if err := requireEmptyWorkspace(ctx, workspace); err != nil {
		return ImportInstalledPluginResult{}, err
	}

	bin, err := a.reader.LoadBinary(ctx, id)
	if err != nil {
		return ImportInstalledPluginResult{}, err
	}

	validated, err := pluginresource.Binary(bin).Validate(ctx)
	if err != nil {
		return ImportInstalledPluginResult{}, err
	}
	if err := validated.Metadata.Validate(); err != nil {
		return ImportInstalledPluginResult{}, err
	}
	if string(record.Manifest.ID) != "" && record.Manifest.ID != validated.Metadata.ID {
		return ImportInstalledPluginResult{}, fmt.Errorf("installed plugin manifest %q does not match archive manifest %q", record.Manifest.ID, validated.Metadata.ID)
	}
	if err := ensurePluginBuildTarget(ctx, string(validated.Metadata.ID), id.String()); err != nil {
		return ImportInstalledPluginResult{}, err
	}

	imported, err := importPluginArchive(ctx, workspace, bin)
	if err != nil {
		return ImportInstalledPluginResult{}, err
	}

	if err := a.setPluginBuildTarget(ctx, pluginBuildTarget{
		Mode:           pluginBuildTargetModeExisting,
		InstallationID: id.String(),
		ManifestID:     string(imported.Manifest.Metadata.ID),
	}); err != nil {
		return ImportInstalledPluginResult{}, err
	}

	return ImportInstalledPluginResult{
		ID:          id.String(),
		Name:        imported.Manifest.Metadata.Name,
		Description: imported.Manifest.Metadata.Description,
		Files:       imported.Files,
	}, nil
}

func requireEmptyWorkspace(ctx context.Context, workspace workspaceprovider.Workspace) error {
	files, err := workspace.List(ctx, workspaceprovider.ListOptions{MaxFiles: 1})
	if err != nil {
		return err
	}
	if len(files) > 0 {
		return errors.New(pluginBuildTargetDifferentPluginMessage)
	}
	return nil
}

func importPluginArchive(ctx context.Context, workspace workspaceprovider.Workspace, bin []byte) (*importedPluginArchive, error) {
	if len(bin) > int(pluginresource.MaxArchiveSizeBytes) {
		return nil, fmt.Errorf("plugin archive exceeds maximum size of %d bytes", pluginresource.MaxArchiveSizeBytes)
	}

	validated, err := pluginresource.Binary(bin).Validate(ctx)
	if err != nil {
		return nil, fmt.Errorf("validate installed plugin archive: %w", err)
	}
	if err := validated.Metadata.Validate(); err != nil {
		return nil, fmt.Errorf("validate installed plugin manifest: %w", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(bin), int64(len(bin)))
	if err != nil {
		return nil, fmt.Errorf("open installed plugin archive: %w", err)
	}
	if len(zr.File) > pluginresource.MaxArchiveExtractFileCount {
		return nil, fmt.Errorf("plugin archive exceeds maximum file count of %d", pluginresource.MaxArchiveExtractFileCount)
	}

	files := []archiveImportFile{}
	var total uint64
	hasRootManifest := false
	for _, file := range zr.File {
		rel, err := cleanPluginArchivePath(file.Name)
		if err != nil {
			return nil, err
		}
		if rel == "" {
			continue
		}
		if rel == pluginresource.ArchiveManifestFileName {
			if file.FileInfo().IsDir() || file.FileInfo().Mode().Type() != 0 {
				return nil, fmt.Errorf("plugin archive contains unsupported file type: %s", rel)
			}
			hasRootManifest = true
			continue
		}
		if file.FileInfo().IsDir() {
			continue
		}
		if file.FileInfo().Mode().Type() != 0 {
			return nil, fmt.Errorf("plugin archive contains unsupported file type: %s", rel)
		}
		if file.UncompressedSize64 > pluginresource.MaxArchiveExtractFileBytes {
			return nil, fmt.Errorf("plugin archive file %s exceeds maximum size of %d bytes", rel, pluginresource.MaxArchiveExtractFileBytes)
		}
		total += file.UncompressedSize64
		if total > pluginresource.MaxArchiveExtractTotalBytes {
			return nil, fmt.Errorf("plugin archive exceeds maximum extracted size of %d bytes", pluginresource.MaxArchiveExtractTotalBytes)
		}

		data, err := readZipFile(file)
		if err != nil {
			return nil, err
		}
		files = append(files, archiveImportFile{Path: rel, Data: data})
	}
	if !hasRootManifest {
		return nil, fmt.Errorf("plugin archive does not contain %s", pluginresource.ArchiveManifestFileName)
	}

	manifestYAML, err := yaml.Marshal(validated.Metadata.ToMap())
	if err != nil {
		return nil, fmt.Errorf("render manifest.yaml: %w", err)
	}
	if _, err := workspace.WriteFile(ctx, manifestYAMLFilename, manifestYAML); err != nil {
		return nil, err
	}

	written := 1
	for _, file := range files {
		if _, err := workspace.WriteFile(ctx, file.Path, file.Data); err != nil {
			return nil, err
		}
		written++
	}

	return &importedPluginArchive{Manifest: validated, Files: written}, nil
}

func cleanPluginArchivePath(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("plugin archive contains an empty file name")
	}
	if strings.Contains(name, "\\") {
		return "", fmt.Errorf("plugin archive path %q is not portable", name)
	}
	for _, part := range strings.Split(name, "/") {
		if part == ".." {
			return "", fmt.Errorf("plugin archive path %q escapes the workspace", name)
		}
	}

	clean := path.Clean(name)
	if clean == "." {
		return "", nil
	}
	if path.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, "../") {
		return "", fmt.Errorf("plugin archive path %q escapes the workspace", name)
	}
	return clean, nil
}

func readZipFile(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	limited := io.LimitReader(rc, int64(file.UncompressedSize64)+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if uint64(len(data)) > file.UncompressedSize64 {
		return nil, fmt.Errorf("plugin archive file %s exceeded its declared size", file.Name)
	}
	if uint64(len(data)) > pluginresource.MaxArchiveExtractFileBytes {
		return nil, fmt.Errorf("plugin archive file %s exceeds maximum size of %d bytes", file.Name, pluginresource.MaxArchiveExtractFileBytes)
	}
	if file.FileInfo().Mode()&fs.ModeSymlink != 0 {
		return nil, fmt.Errorf("plugin archive contains unsupported file type: %s", file.Name)
	}
	return data, nil
}
