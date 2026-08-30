package theme

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_querier"
	"github.com/Southclaws/storyden/app/resources/asset/asset_writer"
	"github.com/Southclaws/storyden/app/resources/audit"
	"github.com/Southclaws/storyden/app/resources/audit/audit_writer"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/admin/settings_manager"
	"github.com/Southclaws/storyden/app/services/asset/asset_download"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/mime"
)

const (
	APIVersion          = "v1"
	MaxAssetSizeBytes   = 1 * 1024 * 1024
	MaxActiveSizeBytes  = 5 * 1024 * 1024
	MaxActiveAssetCount = 32

	metadataPurposeKey   = "purpose"
	metadataPurposeTheme = "theme"
	metadataKindKey      = "theme_kind"
	metadataIntegrityKey = "theme_integrity"
)

var (
	errInvalidAsset = fault.Wrap(fault.New("invalid theme asset"), ftag.With(ftag.InvalidArgument))
	errInvalidTheme = fault.Wrap(fault.New("invalid theme manifest"), ftag.With(ftag.InvalidArgument))
	errActiveAsset  = fault.Wrap(fault.New("active theme assets cannot be deleted"), ftag.With(ftag.InvalidArgument))
)

type AssetKind string

const (
	AssetKindStylesheet AssetKind = "stylesheet"
	AssetKindScript     AssetKind = "script"
)

func (k AssetKind) MIME() string {
	if k == AssetKindStylesheet {
		return "text/css"
	}
	return "application/javascript"
}

type Asset struct {
	ID        string
	Filename  string
	Path      string
	MIMEType  string
	Size      int
	Integrity string
	Kind      AssetKind
}

type Manifest struct {
	APIVersion      string
	Enabled         bool
	Revision        string
	Stylesheets     []Asset
	Scripts         []Asset
	RuntimeDisabled bool
}

type Service struct {
	settings   *settings_manager.Manager
	assets     *asset_writer.Writer
	assetQuery *asset_querier.Querier
	downloader *asset_download.Downloader
	objects    object.Storer
	audit      *audit_writer.Writer
	disabled   bool
}

func New(
	cfg config.Config,
	settings *settings_manager.Manager,
	assets *asset_writer.Writer,
	assetQuery *asset_querier.Querier,
	downloader *asset_download.Downloader,
	objects object.Storer,
	auditWriter *audit_writer.Writer,
) *Service {
	return &Service{
		settings:   settings,
		assets:     assets,
		assetQuery: assetQuery,
		downloader: downloader,
		objects:    objects,
		audit:      auditWriter,
		disabled:   cfg.CustomThemesDisable,
	}
}

func EmptyManifest() Manifest {
	return Manifest{
		APIVersion:  APIVersion,
		Stylesheets: []Asset{},
		Scripts:     []Asset{},
	}
}

func (s *Service) GetPublic(ctx context.Context) (Manifest, error) {
	manifest, err := s.getConfigured(ctx)
	if err != nil {
		return EmptyManifest(), err
	}
	if s.disabled {
		out := EmptyManifest()
		out.RuntimeDisabled = true
		return out, nil
	}
	return manifest, nil
}

func (s *Service) GetAdmin(ctx context.Context) (Manifest, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return EmptyManifest(), fault.Wrap(err, fctx.With(ctx))
	}
	manifest, err := s.getConfigured(ctx)
	if err != nil {
		return EmptyManifest(), err
	}
	manifest.RuntimeDisabled = s.disabled
	return manifest, nil
}

func (s *Service) getConfigured(ctx context.Context) (Manifest, error) {
	set, err := s.settings.Get(ctx)
	if err != nil {
		return EmptyManifest(), fault.Wrap(err, fctx.With(ctx))
	}

	themeSettings := set.Theme.OrZero()
	stylesheets := mapStoredAssets(themeSettings.Stylesheets, AssetKindStylesheet)
	scripts := mapStoredAssets(themeSettings.Scripts, AssetKindScript)
	revision := revisionFor(stylesheets, scripts)

	return Manifest{
		APIVersion:  APIVersion,
		Enabled:     len(stylesheets)+len(scripts) > 0,
		Revision:    revision,
		Stylesheets: stylesheets,
		Scripts:     scripts,
	}, nil
}

func (s *Service) Upload(ctx context.Context, r io.Reader, size int64, clientFilename string) (Asset, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return Asset{}, fault.Wrap(err, fctx.With(ctx))
	}
	if size < 1 || size > MaxAssetSizeBytes {
		return Asset{}, fault.Wrap(errInvalidAsset, fmsg.Withf("theme assets must be between 1 and %d bytes", MaxAssetSizeBytes))
	}

	kind, err := classify(clientFilename)
	if err != nil {
		return Asset{}, err
	}

	content, err := io.ReadAll(io.LimitReader(r, MaxAssetSizeBytes+1))
	if err != nil {
		return Asset{}, fault.Wrap(err, fctx.With(ctx))
	}
	if int64(len(content)) != size || len(content) > MaxAssetSizeBytes || !utf8.Valid(content) || bytes.IndexByte(content, 0) >= 0 {
		return Asset{}, fault.Wrap(errInvalidAsset, fmsg.With("theme assets must be valid UTF-8 text matching Content-Length"))
	}

	hash := sha256.Sum256(content)
	integrity := "sha256-" + base64.StdEncoding.EncodeToString(hash[:])
	filename := asset.NewFilename(clientFilename)
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return Asset{}, fault.Wrap(err, fctx.With(ctx))
	}

	path := asset.BuildAssetPath(filename)
	if err := s.objects.Write(ctx, path, bytes.NewReader(content), int64(len(content))); err != nil {
		return Asset{}, fault.Wrap(err, fctx.With(ctx))
	}

	stored, err := s.assets.AddWithMetadata(
		ctx,
		xid.ID(accountID),
		filename,
		len(content),
		mime.New(kind.MIME()),
		asset.Metadata{
			metadataPurposeKey:   metadataPurposeTheme,
			metadataKindKey:      string(kind),
			metadataIntegrityKey: integrity,
		},
	)
	if err != nil {
		_ = s.objects.Delete(ctx, path)
		return Asset{}, fault.Wrap(err, fctx.With(ctx))
	}

	if manifest, getErr := s.getConfigured(ctx); getErr == nil {
		s.cleanup(ctx, manifest)
	}
	return mapAsset(stored, kind), nil
}

func (s *Service) Publish(ctx context.Context, stylesheetIDs, scriptIDs []string) (Manifest, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return EmptyManifest(), fault.Wrap(err, fctx.With(ctx))
	}
	if len(stylesheetIDs)+len(scriptIDs) > MaxActiveAssetCount {
		return EmptyManifest(), fault.Wrap(errInvalidTheme, fmsg.Withf("a theme may contain at most %d active assets", MaxActiveAssetCount))
	}

	stylesheets, err := s.resolve(ctx, stylesheetIDs, AssetKindStylesheet)
	if err != nil {
		return EmptyManifest(), err
	}
	scripts, err := s.resolve(ctx, scriptIDs, AssetKindScript)
	if err != nil {
		return EmptyManifest(), err
	}

	total := 0
	for _, item := range append(append([]Asset{}, stylesheets...), scripts...) {
		total += item.Size
	}
	if total > MaxActiveSizeBytes {
		return EmptyManifest(), fault.Wrap(errInvalidTheme, fmsg.Withf("active theme assets may total at most %d bytes", MaxActiveSizeBytes))
	}

	_, err = s.settings.Set(ctx, settings.Settings{
		Theme: opt.New(settings.ThemeSettings{
			Stylesheets: mapThemeSettingsAssets(stylesheets),
			Scripts:     mapThemeSettingsAssets(scripts),
		}),
	})
	if err != nil {
		return EmptyManifest(), fault.Wrap(err, fctx.With(ctx))
	}

	manifest := Manifest{
		APIVersion:      APIVersion,
		Enabled:         len(stylesheets)+len(scripts) > 0,
		Revision:        revisionFor(stylesheets, scripts),
		Stylesheets:     stylesheets,
		Scripts:         scripts,
		RuntimeDisabled: s.disabled,
	}
	if err := s.recordAudit(ctx, audit.EventTypeThemePublished, manifest); err != nil {
		return EmptyManifest(), err
	}
	s.cleanup(ctx, manifest)
	return manifest, nil
}

func (s *Service) Disable(ctx context.Context) error {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	previous, err := s.getConfigured(ctx)
	if err != nil {
		return err
	}
	if _, err := s.settings.Set(ctx, settings.Settings{Theme: opt.New(settings.ThemeSettings{Stylesheets: []settings.ThemeAsset{}, Scripts: []settings.ThemeAsset{}})}); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	return s.recordAudit(ctx, audit.EventTypeThemeDisabled, previous)
}

func (s *Service) recordAudit(ctx context.Context, eventType audit.EventType, manifest Manifest) error {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	assets := make([]map[string]any, 0, len(manifest.Stylesheets)+len(manifest.Scripts))
	for _, item := range append(append([]Asset{}, manifest.Stylesheets...), manifest.Scripts...) {
		assets = append(assets, map[string]any{"filename": item.Filename, "kind": string(item.Kind), "size": item.Size, "integrity": item.Integrity})
	}
	_, err = s.audit.Create(ctx, eventType, opt.New(account.AccountID(accountID)), opt.NewEmpty[datagraph.Ref](), map[string]any{"revision": manifest.Revision, "assets": assets})
	return fault.Wrap(err, fctx.With(ctx))
}

func (s *Service) GetAsset(ctx context.Context, filename string) (Asset, io.Reader, error) {
	a, r, err := s.downloader.Get(ctx, asset.NewFilepathFilename(filename))
	if err != nil {
		return Asset{}, nil, fault.Wrap(err, fctx.With(ctx))
	}
	kind, ok := themeAssetKind(a)
	if !ok {
		return Asset{}, nil, fault.Wrap(errInvalidAsset, ftag.With(ftag.NotFound))
	}
	return mapAsset(a, kind), r, nil
}

func (s *Service) DeleteAsset(ctx context.Context, filename string) error {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	a, err := s.assetQuery.Get(ctx, asset.NewFilepathFilename(filename))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if _, ok := themeAssetKind(a); !ok {
		return fault.Wrap(errInvalidAsset, ftag.With(ftag.NotFound))
	}

	manifest, err := s.getConfigured(ctx)
	if err != nil {
		return err
	}
	for _, item := range append(append([]Asset{}, manifest.Stylesheets...), manifest.Scripts...) {
		if item.Filename == a.Name.String() {
			return errActiveAsset
		}
	}

	if err := s.objects.Delete(ctx, asset.BuildAssetPath(a.Name)); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if err := s.assets.Remove(ctx, xid.ID(accountID), a.Name); err != nil {
		return err
	}
	s.cleanup(ctx, manifest)
	return nil
}

// cleanup opportunistically removes abandoned uploads. It is deliberately
// best-effort so storage maintenance can never make publication unavailable.
func (s *Service) cleanup(ctx context.Context, manifest Manifest) {
	rows, err := s.assetQuery.ListOlderThan(ctx, time.Now().Add(-24*time.Hour))
	if err != nil {
		return
	}
	active := map[string]struct{}{}
	for _, item := range append(append([]Asset{}, manifest.Stylesheets...), manifest.Scripts...) {
		active[item.ID] = struct{}{}
	}
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return
	}
	for _, row := range rows {
		if _, ok := themeAssetKind(row); !ok {
			continue
		}
		if _, ok := active[row.ID.String()]; ok {
			continue
		}
		if err := s.objects.Delete(ctx, asset.BuildAssetPath(row.Name)); err != nil {
			continue
		}
		_ = s.assets.Remove(ctx, xid.ID(accountID), row.Name)
	}
}

func (s *Service) resolve(ctx context.Context, ids []string, expected AssetKind) ([]Asset, error) {
	seen := map[string]struct{}{}
	out := make([]Asset, 0, len(ids))
	for _, raw := range ids {
		if _, ok := seen[raw]; ok {
			continue
		}
		id, err := xid.FromString(raw)
		if err != nil {
			return nil, fault.Wrap(errInvalidTheme, fmsg.Withf("invalid asset id %q", raw))
		}
		a, err := s.assetQuery.GetByID(ctx, id)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		kind, ok := themeAssetKind(a)
		if !ok || kind != expected {
			return nil, fault.Wrap(errInvalidTheme, fmsg.Withf("asset %q is not a theme %s", raw, expected))
		}
		seen[raw] = struct{}{}
		out = append(out, mapAsset(a, kind))
	}
	return out, nil
}

func classify(name string) (AssetKind, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" || strings.ContainsAny(trimmed, `/\\`) || filepath.Base(trimmed) != trimmed {
		return "", fault.Wrap(errInvalidAsset, fmsg.With("filename must be a plain .css or .js filename"))
	}
	switch strings.ToLower(filepath.Ext(trimmed)) {
	case ".css":
		return AssetKindStylesheet, nil
	case ".js":
		return AssetKindScript, nil
	default:
		return "", fault.Wrap(errInvalidAsset, fmsg.With("theme assets must use .css or .js"))
	}
}

func themeAssetKind(a *asset.Asset) (AssetKind, bool) {
	if a.Metadata[metadataPurposeKey] != metadataPurposeTheme {
		return "", false
	}
	kind, ok := a.Metadata[metadataKindKey].(string)
	if !ok || (AssetKind(kind) != AssetKindStylesheet && AssetKind(kind) != AssetKindScript) {
		return "", false
	}
	return AssetKind(kind), true
}

func mapAsset(a *asset.Asset, kind AssetKind) Asset {
	integrity, _ := a.Metadata[metadataIntegrityKey].(string)
	return Asset{
		ID:        a.ID.String(),
		Filename:  a.Name.String(),
		Path:      fmt.Sprintf("/api/info/theme/assets/%s", a.Name.String()),
		MIMEType:  kind.MIME(),
		Size:      a.Size,
		Integrity: integrity,
		Kind:      kind,
	}
}

func mapStoredAssets(items []settings.ThemeAsset, kind AssetKind) []Asset {
	out := make([]Asset, 0, len(items))
	for _, item := range items {
		out = append(out, Asset{
			ID:        item.ID,
			Filename:  item.Filename,
			Path:      fmt.Sprintf("/api/info/theme/assets/%s", item.Filename),
			MIMEType:  item.MIMEType,
			Size:      item.Size,
			Integrity: item.Integrity,
			Kind:      kind,
		})
	}
	return out
}

func mapThemeSettingsAssets(items []Asset) []settings.ThemeAsset {
	out := make([]settings.ThemeAsset, 0, len(items))
	for _, item := range items {
		out = append(out, settings.ThemeAsset{
			ID:        item.ID,
			Filename:  item.Filename,
			MIMEType:  item.MIMEType,
			Size:      item.Size,
			Integrity: item.Integrity,
		})
	}
	return out
}

func revisionFor(stylesheets, scripts []Asset) string {
	if len(stylesheets)+len(scripts) == 0 {
		return ""
	}
	h := sha256.New()
	for _, item := range stylesheets {
		_, _ = io.WriteString(h, "css\x00"+item.ID+"\x00")
	}
	for _, item := range scripts {
		_, _ = io.WriteString(h, "js\x00"+item.ID+"\x00")
	}
	return hex.EncodeToString(h.Sum(nil))
}
