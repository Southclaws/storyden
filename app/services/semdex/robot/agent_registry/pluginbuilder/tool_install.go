package pluginbuilder

import (
	"bytes"
	"context"
	"fmt"

	"github.com/rs/xid"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	pluginresource "github.com/Southclaws/storyden/app/resources/plugin"
)

type InstallInput struct {
	InstallationID string `json:"installation_id,omitempty" jsonschema:"Existing supervised plugin installation ID to update"`
	Activate       bool   `json:"activate,omitempty" jsonschema:"Activate or restart the supervised plugin after installing"`
	UpdateIfExists bool   `json:"update_if_exists,omitempty" jsonschema:"Update the first installed supervised plugin with the same manifest ID"`
	RunValidation  bool   `json:"run_validation,omitempty" jsonschema:"Deprecated; validation runs by default unless skip_validation is true"`
	SkipValidation bool   `json:"skip_validation,omitempty" jsonschema:"Skip pre-install validation commands"`
}

type InstallResult struct {
	Action         string `json:"action"`
	InstallationID string `json:"installation_id"`
	ManifestID     string `json:"manifest_id"`
	Active         bool   `json:"active"`
	PackageBytes   int    `json:"package_bytes"`
}

func (a *Agent) addInstallTools(add toolAdder) error {
	return add(functiontool.New(functiontool.Config{
		Name:        "plugin_install",
		Description: "Package and install or update a managed plugin as a supervised Storyden plugin. Requires administrator context.",
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
		if result, err := a.GoFormat(ctx); err != nil {
			return InstallResult{}, err
		} else if !result.Success {
			return InstallResult{}, fmt.Errorf("gofmt failed: %s", result.Output)
		}
		if result, err := a.GoTidy(ctx); err != nil {
			return InstallResult{}, err
		} else if !result.Success {
			return InstallResult{}, fmt.Errorf("go mod tidy failed: %s", result.Output)
		}
		if result, err := a.GoVet(ctx); err != nil {
			return InstallResult{}, err
		} else if !result.Success {
			return InstallResult{}, fmt.Errorf("go vet failed: %s", result.Output)
		}
		if result, err := a.GoTest(ctx, GoTestInput{}); err != nil {
			return InstallResult{}, err
		} else if !result.Success {
			return InstallResult{}, fmt.Errorf("go test failed: %s", result.Output)
		}
	}

	pkg, err := a.PackageBytes(ctx)
	if err != nil {
		return InstallResult{}, err
	}

	id, found, action, err := a.resolveInstallation(ctx, in.InstallationID, string(pkg.Manifest.ID), in.UpdateIfExists)
	if err != nil {
		return InstallResult{}, err
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

func (a *Agent) resolveInstallation(ctx context.Context, installationID, manifestID string, updateIfExists bool) (pluginresource.InstallationID, bool, string, error) {
	if installationID != "" {
		id, err := xid.FromString(installationID)
		if err != nil {
			return pluginresource.InstallationID{}, false, "", fmt.Errorf("invalid installation_id: %w", err)
		}
		return pluginresource.InstallationID(id), true, "updated", nil
	}

	if !updateIfExists {
		return pluginresource.InstallationID{}, false, "", nil
	}

	records, err := a.installer.List(ctx)
	if err != nil {
		return pluginresource.InstallationID{}, false, "", err
	}

	for _, rec := range records {
		if rec.Mode.Supervised() && string(rec.Manifest.ID) == manifestID {
			return rec.InstallationID, true, "updated", nil
		}
	}

	return pluginresource.InstallationID{}, false, "", nil
}
