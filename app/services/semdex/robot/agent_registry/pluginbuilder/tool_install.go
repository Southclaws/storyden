package pluginbuilder

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/rs/xid"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	pluginresource "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type InstallInput struct {
	Activate       bool `json:"activate,omitempty" jsonschema:"Activate or restart the supervised plugin after installing"`
	SkipValidation bool `json:"skip_validation,omitempty" jsonschema:"Skip pre-install validation commands"`
}

type InstallResult struct {
	Action         string `json:"action"`
	InstallationID string `json:"installation_id"`
	ManifestID     string `json:"manifest_id"`
	Active         bool   `json:"active"`
	PackageBytes   int    `json:"package_bytes"`
}

type packageArchive struct {
	Manifest rpc.Manifest
	Bytes    []byte
	Files    []string
}

type manifestFile struct {
	Path     string
	Manifest rpc.Manifest
}

type pluginBuildRuntimeTarget struct {
	GOOS   string
	GOARCH string
}

const (
	packagedPluginBinary  = "main.exe"
	packagedPluginCommand = "./" + packagedPluginBinary
	pluginBuildTimeout    = 5 * time.Minute
	pluginBuildFileWait   = 10 * time.Second
)

func (a *Agent) addInstallTools(add toolAdder) error {
	return add(functiontool.New(functiontool.Config{
		Name:        "plugin_install",
		Description: "Package and install or update a managed plugin as a supervised Storyden plugin.",
	}, func(ctx adktool.Context, args InstallInput) (InstallResult, error) {
		result, err := a.Install(ctx, args)
		if err != nil {
			return InstallResult{}, err
		}
		return result, nil
	}))
}

func (a *Agent) Install(ctx context.Context, in InstallInput) (InstallResult, error) {
	if a.installer == nil {
		return InstallResult{}, fmt.Errorf("plugin installer is not configured")
	}

	if !in.SkipValidation {
		if result, err := a.Validate(ctx, ValidateInput{}); err != nil {
			return InstallResult{}, err
		} else if !result.Success {
			return InstallResult{}, errors.New(result.Message)
		}
	}

	pkg, err := a.packageBytes(ctx)
	if err != nil {
		return InstallResult{}, err
	}

	id, found, action, err := a.resolveInstallation(ctx, string(pkg.Manifest.ID))
	if err != nil {
		return InstallResult{}, err
	}
	targetMode := pluginBuildTargetModeNew
	if target, ok, err := pluginBuildTargetFromContext(ctx); err != nil {
		return InstallResult{}, err
	} else if ok {
		if target.Mode != "" {
			targetMode = target.Mode
		}
	}

	active := false
	if !found {
		available, err := a.installer.AddFromFile(ctx, bytes.NewReader(pkg.Bytes))
		if err != nil {
			return InstallResult{}, err
		}
		id = available.Record.InstallationID
		action = "installed"
	} else {
		rec, err := a.installer.UpdatePackage(ctx, id, bytes.NewReader(pkg.Bytes))
		if err != nil {
			return InstallResult{}, err
		}
		id = rec.InstallationID
		if action == "" {
			action = "updated"
		}
	}

	if err := a.setPluginBuildTarget(ctx, pluginBuildTarget{
		Mode:           targetMode,
		InstallationID: id.String(),
		ManifestID:     string(pkg.Manifest.ID),
	}); err != nil {
		return InstallResult{}, err
	}

	if in.Activate {
		if err := a.installer.SetActiveState(ctx, id, pluginresource.ActiveStateActive); err != nil {
			return InstallResult{}, err
		}
		active = true
	}

	return InstallResult{
		Action:         action,
		InstallationID: id.String(),
		ManifestID:     string(pkg.Manifest.ID),
		Active:         active,
		PackageBytes:   len(pkg.Bytes),
	}, nil
}

func (a *Agent) resolveInstallation(ctx context.Context, manifestID string) (pluginresource.InstallationID, bool, string, error) {
	target, ok, err := pluginBuildTargetFromContext(ctx)
	if err != nil {
		return pluginresource.InstallationID{}, false, "", err
	}
	if ok {
		if target.ManifestID != "" && target.ManifestID != manifestID {
			return pluginresource.InstallationID{}, false, "", errors.New(pluginBuildTargetDifferentPluginMessage)
		}
		if target.InstallationID == "" {
			if target.Mode == pluginBuildTargetModeNew {
				return pluginresource.InstallationID{}, false, "", nil
			}
			return pluginresource.InstallationID{}, false, "", errors.New("this chat cannot update the selected plugin; start a new chat for that plugin")
		}

		id, err := xid.FromString(target.InstallationID)
		if err != nil {
			return pluginresource.InstallationID{}, false, "", fmt.Errorf("invalid bound plugin installation: %w", err)
		}
		return pluginresource.InstallationID(id), true, "updated", nil
	}

	if id, ok, err := a.findExistingSupervisedInstallation(ctx, manifestID); err != nil || ok {
		return id, ok, "updated", err
	}

	return pluginresource.InstallationID{}, false, "", nil
}

func (a *Agent) findExistingSupervisedInstallation(ctx context.Context, manifestID string) (pluginresource.InstallationID, bool, error) {
	if a == nil || a.installer == nil || strings.TrimSpace(manifestID) == "" {
		return pluginresource.InstallationID{}, false, nil
	}

	records, err := a.installer.List(ctx)
	if err != nil {
		return pluginresource.InstallationID{}, false, err
	}

	var newest *pluginresource.Record
	for _, rec := range records {
		if rec == nil || rec.Manifest.ID != manifestID {
			continue
		}
		if rec.Mode.External() {
			return pluginresource.InstallationID{}, false, fmt.Errorf("plugin manifest id %q is already installed as an external plugin", manifestID)
		}
		if !rec.Mode.Supervised() {
			continue
		}
		if newest == nil || rec.Created.After(newest.Created) {
			newest = rec
		}
	}
	if newest == nil {
		return pluginresource.InstallationID{}, false, nil
	}

	return newest.InstallationID, true, nil
}

func (a *Agent) packageBytes(ctx context.Context) (*packageArchive, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return nil, err
	}
	return buildPackage(ctx, workspace)
}

func buildPackage(ctx context.Context, workspace workspaceprovider.Workspace) (*packageArchive, error) {
	return buildPackageForTarget(ctx, workspace, hostPluginBuildRuntimeTarget())
}

func buildPackageForTarget(ctx context.Context, workspace workspaceprovider.Workspace, target pluginBuildRuntimeTarget) (*packageArchive, error) {
	mf, err := readProjectManifest(ctx, workspace)
	if err != nil {
		return nil, err
	}

	sourceFiles, err := packageWorkspaceFiles(ctx, workspace)
	if err != nil {
		return nil, err
	}
	if err := validateHostAPIAccessManifest(ctx, workspace, mf.Manifest, sourceFiles); err != nil {
		return nil, err
	}

	if err := buildPackagedPluginBinary(ctx, workspace, target); err != nil {
		return nil, err
	}

	packagedManifest := manifestWithPackagedRuntime(mf.Manifest)
	manifestJSON, err := json.MarshalIndent(packagedManifest, "", "  ")
	if err != nil {
		return nil, err
	}
	manifestJSON = append(manifestJSON, '\n')

	binaryFile, err := waitForWorkspaceReadFile(ctx, workspace, packagedPluginBinary, pluginBuildFileWait)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	manifestHeader := &zip.FileHeader{Name: pluginresource.ArchiveManifestFileName, Method: zip.Deflate}
	manifestHeader.SetMode(0o644)
	writer, err := zw.CreateHeader(manifestHeader)
	if err != nil {
		return nil, err
	}
	if _, err := writer.Write(manifestJSON); err != nil {
		return nil, err
	}

	written := []string{pluginresource.ArchiveManifestFileName}
	for _, file := range sourceFiles {
		if file.Path == packagedPluginBinary {
			continue
		}
		data, err := workspace.ReadFile(ctx, file.Path, -1)
		if err != nil {
			return nil, err
		}
		if err := writeArchiveFile(zw, file.Path, file.Mode, data.Content); err != nil {
			return nil, err
		}
		written = append(written, file.Path)
	}
	if err := writeArchiveFile(zw, packagedPluginBinary, 0o755, binaryFile.Content); err != nil {
		return nil, err
	}
	written = append(written, packagedPluginBinary)

	if err := zw.Close(); err != nil {
		return nil, err
	}

	validated, err := pluginresource.Binary(buf.Bytes()).Validate(ctx)
	if err != nil {
		return nil, fmt.Errorf("validate package: %w", err)
	}
	if err := validated.Metadata.Validate(); err != nil {
		return nil, fmt.Errorf("validate package manifest: %w", err)
	}

	return &packageArchive{Manifest: packagedManifest, Bytes: buf.Bytes(), Files: written}, nil
}

func writeArchiveFile(zw *zip.Writer, path string, mode fs.FileMode, content []byte) error {
	header := &zip.FileHeader{Name: path, Method: zip.Deflate}
	header.SetMode(mode)
	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	if _, err := writer.Write(content); err != nil {
		return err
	}
	return nil
}

func waitForWorkspaceReadFile(ctx context.Context, workspace workspaceprovider.Workspace, path string, timeout time.Duration) (workspaceprovider.ReadFileResult, error) {
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	var lastErr error
	for {
		data, err := workspace.ReadFile(ctx, path, -1)
		if err == nil {
			return data, nil
		} else {
			lastErr = err
		}

		select {
		case <-ctx.Done():
			if lastErr != nil {
				return workspaceprovider.ReadFileResult{}, lastErr
			}
			return workspaceprovider.ReadFileResult{}, ctx.Err()
		case <-deadline.C:
			if lastErr != nil {
				return workspaceprovider.ReadFileResult{}, fmt.Errorf("build plugin binary: %s was not readable after go build: %w", path, lastErr)
			}
			return workspaceprovider.ReadFileResult{}, fmt.Errorf("build plugin binary: %s was not readable after go build", path)
		case <-ticker.C:
		}
	}
}

func hostPluginBuildRuntimeTarget() pluginBuildRuntimeTarget {
	return pluginBuildRuntimeTarget{
		GOOS:   runtime.GOOS,
		GOARCH: runtime.GOARCH,
	}
}

func buildPackagedPluginBinary(ctx context.Context, workspace workspaceprovider.Workspace, target pluginBuildRuntimeTarget) error {
	if target.GOOS == "" || target.GOARCH == "" {
		return errors.New("plugin build target GOOS and GOARCH are required")
	}

	result, err := workspace.Run(ctx, workspaceprovider.CommandSpec{
		Command: "go",
		Args:    []string{"build", "-trimpath", "-o", packagedPluginBinary, "."},
		Env: []string{
			"GOOS=" + target.GOOS,
			"GOARCH=" + target.GOARCH,
			"CGO_ENABLED=0",
		},
		Timeout: pluginBuildTimeout,
	})
	if err != nil {
		return fmt.Errorf("build plugin binary for %s/%s: %w", target.GOOS, target.GOARCH, err)
	}
	if !result.Success {
		message := strings.TrimSpace(result.Output)
		if message == "" {
			message = strings.TrimSpace(result.Error)
		}
		if message == "" {
			message = "go build failed"
		}
		return fmt.Errorf("build plugin binary for %s/%s: %s", target.GOOS, target.GOARCH, message)
	}

	return nil
}

func manifestWithPackagedRuntime(manifest rpc.Manifest) rpc.Manifest {
	manifest.Command = packagedPluginCommand
	manifest.Args = nil
	return manifest
}

func packageWorkspaceFiles(ctx context.Context, workspace workspaceprovider.Workspace) ([]workspaceprovider.FileInfo, error) {
	listed, err := workspace.List(ctx, workspaceprovider.ListOptions{MaxFiles: 500})
	if err != nil {
		return nil, err
	}

	files := []workspaceprovider.FileInfo{}
	for _, file := range listed {
		name := file.Path
		if name == ".DS_Store" || name == pluginresource.ArchiveManifestFileName {
			continue
		}
		if strings.HasSuffix(name, ".zip") || name == manifestYAMLFilename {
			continue
		}
		if file.Mode&fs.ModeSymlink != 0 {
			continue
		}
		files = append(files, file)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })

	return files, nil
}

func readProjectManifest(ctx context.Context, workspace workspaceprovider.Workspace) (*manifestFile, error) {
	data, err := workspace.ReadFile(ctx, manifestYAMLFilename, -1)
	if err != nil {
		return nil, err
	}
	var raw map[string]any
	if err := yaml.Unmarshal(data.Content, &raw); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	if err := validateManifestRaw(raw); err != nil {
		return nil, fmt.Errorf("validate manifest: %w", err)
	}
	manifest, err := rpc.ManifestFromMap(raw)
	if err != nil {
		return nil, fmt.Errorf("validate manifest: %w", err)
	}
	return &manifestFile{Path: data.Path, Manifest: *manifest}, nil
}

func validateHostAPIAccessManifest(ctx context.Context, workspace workspaceprovider.Workspace, manifest rpc.Manifest, files []workspaceprovider.FileInfo) error {
	usesHostAPI, err := workspaceUsesBuildAPIClient(ctx, workspace, files)
	if err != nil {
		return err
	}
	if !usesHostAPI {
		return nil
	}
	if _, ok := manifest.Access.Get(); ok {
		return nil
	}

	return errors.New("manifest access is required because plugin code uses BuildAPIClient; add access with a stable bot account handle, display name, and narrow Storyden permissions for the API operations being called")
}

func workspaceUsesBuildAPIClient(ctx context.Context, workspace workspaceprovider.Workspace, files []workspaceprovider.FileInfo) (bool, error) {
	for _, file := range files {
		if !strings.HasSuffix(file.Path, ".go") {
			continue
		}
		data, err := workspace.ReadFile(ctx, file.Path, -1)
		if err != nil {
			return false, err
		}
		if bytes.Contains(data.Content, []byte("BuildAPIClient(")) {
			return true, nil
		}
	}

	return false, nil
}
