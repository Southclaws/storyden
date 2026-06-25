package pluginbuilder

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	pluginresource "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type PackageResult struct {
	ManifestID string   `json:"manifest_id"`
	Path       string   `json:"path"`
	Bytes      int      `json:"bytes"`
	Files      []string `json:"files"`
}

type PackageArchive struct {
	Manifest rpc.Manifest
	Bytes    []byte
	Files    []string
}

type manifestFile struct {
	Path     string
	Manifest rpc.Manifest
}

func (a *Agent) addPackageTools(add toolAdder) error {
	return add(functiontool.New(functiontool.Config{
		Name:        "plugin_package",
		Description: "Build and validate a supervised Storyden plugin package zip from the managed workspace.",
	}, func(ctx adktool.Context, args struct{}) (PackageResult, error) {
		result, err := a.Package(ctx)
		if err != nil {
			return PackageResult{}, err
		}
		return result, nil
	}))
}

func (a *Agent) Package(ctx context.Context) (PackageResult, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return PackageResult{}, err
	}

	pkg, err := buildPackage(ctx, workspace)
	if err != nil {
		return PackageResult{}, err
	}

	path := slugify(string(pkg.Manifest.ID)) + ".zip"
	written, err := workspace.WriteFile(ctx, path, pkg.Bytes)
	if err != nil {
		return PackageResult{}, err
	}
	return PackageResult{
		ManifestID: string(pkg.Manifest.ID),
		Path:       written.Path,
		Bytes:      len(pkg.Bytes),
		Files:      pkg.Files,
	}, nil
}

func (a *Agent) PackageBytes(ctx context.Context) (*PackageArchive, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return nil, err
	}
	return buildPackage(ctx, workspace)
}

func buildPackage(ctx context.Context, workspace workspaceprovider.Workspace) (*PackageArchive, error) {
	mf, err := readProjectManifest(ctx, workspace)
	if err != nil {
		return nil, err
	}

	manifestJSON, err := json.MarshalIndent(mf.Manifest, "", "  ")
	if err != nil {
		return nil, err
	}
	manifestJSON = append(manifestJSON, '\n')

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
	for _, file := range files {
		data, err := workspace.ReadFile(ctx, file.Path, -1)
		if err != nil {
			return nil, err
		}
		header := &zip.FileHeader{Name: file.Path, Method: zip.Deflate}
		header.SetMode(file.Mode)
		header.Method = zip.Deflate
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return nil, err
		}
		if _, err := writer.Write(data.Content); err != nil {
			return nil, err
		}
		written = append(written, file.Path)
	}

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

	return &PackageArchive{Manifest: mf.Manifest, Bytes: buf.Bytes(), Files: written}, nil
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
	manifest, err := rpc.ManifestFromMap(raw)
	if err != nil {
		return nil, fmt.Errorf("validate manifest: %w", err)
	}
	return &manifestFile{Path: data.Path, Manifest: *manifest}, nil
}
