package plugin

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"path/filepath"
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

const PluginDirectory = "plugins"

type InstallationID xid.ID

func (id InstallationID) String() string {
	return xid.ID(id).String()
}

// Record represents the database record for a plugin. It's the source of truth
// for plugins in general, it stores a copy of the manifest, who added it and
// the file path for the plugin's packaged binary application.
type Record struct {
	InstallationID InstallationID
	Created        time.Time
	AddedBy        account.AccountID
	Manifest       rpc.Manifest
	Config         map[string]any
	Token          string
	Mode           PluginMode

	// Supervised tells the host whether the plugin is a process owned by the
	// process. When it is, the host will start the process and monitor it. If
	// false, the plugin is an "external" plugin and not run as a child process.
	Supervised bool

	State          ActiveState
	StateChangedAt time.Time
	StatusMessage  string
	Details        map[string]any

	ReportedState ReportedState
	StartedAt     time.Time
}

type Loaded struct {
	Record Record
	State  ActiveState
}

// Validated represents a plugin that has been read and validated from its binary.
type Validated struct {
	Metadata rpc.Manifest
	Binary   []byte
}

// Available represents a plugin that has been installed and registered in the
// database and is ready to use.
type Available struct {
	Record Record
	Loaded *Validated
}

type Binary []byte

const ArchiveManifestFileName = "manifest.json"

func (b Binary) Validate(ctx context.Context) (*Validated, error) {
	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to open archive"))
	}

	var data []byte
	for _, file := range zr.File {
		if filepath.Base(file.Name) != ArchiveManifestFileName {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return nil, fault.Wrap(err, fmsg.With("failed to open manifest"))
		}
		defer rc.Close()

		data, err = io.ReadAll(rc)
		if err != nil {
			return nil, fault.Wrap(err, fmsg.With("failed to read manifest"))
		}
		break
	}
	if data == nil {
		return nil, fault.New("manifest not found in archive")
	}

	var m rpc.Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fault.Wrap(err, fmsg.With("plugin manifest is malformed"))
	}

	return &Validated{
		Metadata: m,
		Binary:   b,
	}, nil
}

func MapRecord(in *ent.Plugin) (*Record, error) {
	desiredActiveState, err := NewActiveState(in.ActiveState)
	if err != nil {
		return nil, err
	}

	manifest, err := rpc.ManifestFromMap(in.Manifest)
	if err != nil {
		return nil, err
	}

	return &Record{
		InstallationID: InstallationID(in.ID),
		Created:        in.CreatedAt,
		AddedBy:        account.AccountID(in.AddedBy),
		Manifest:       *manifest,
		Config:         mapConfig(in.Config),
		Mode:           pluginModeFromSupervised(in.Supervised),
		Supervised:     in.Supervised, // legacy field; prefer Mode
		Token:          mapToken(pluginModeFromSupervised(in.Supervised), in.AuthSecret),

		State:          desiredActiveState,
		StateChangedAt: in.ActiveStateChangedAt,
		StatusMessage:  opt.NewPtr(in.StatusMessage).OrZero(),
		Details:        in.StatusDetails,
	}, nil
}

func mapToken(mode PluginMode, authSecret string) string {
	if mode.Supervised() {
		return ""
	}
	return authSecret
}

func mapConfig(in map[string]any) map[string]any {
	if in == nil {
		return map[string]any{}
	}
	return in
}
