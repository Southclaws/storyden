package theme

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_querier"
	"github.com/Southclaws/storyden/app/resources/asset/asset_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/admin/settings_manager"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/mime"
)

const (
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

type Service struct {
	settings   *settings_manager.Manager
	assets     *asset_writer.Writer
	assetQuery *asset_querier.Querier
	objects    object.Storer
}

func New(
	settings *settings_manager.Manager,
	assets *asset_writer.Writer,
	assetQuery *asset_querier.Querier,
	objects object.Storer,
) *Service {
	return &Service{
		settings:   settings,
		assets:     assets,
		assetQuery: assetQuery,
		objects:    objects,
	}
}

func (s *Service) GetPublic(ctx context.Context) (settings.ThemeSettings, error) {
	return s.getConfigured(ctx)
}

func (s *Service) GetAdmin(ctx context.Context) (settings.ThemeSettings, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return settings.ThemeSettings{}, fault.Wrap(err, fctx.With(ctx))
	}
	return s.getConfigured(ctx)
}

func (s *Service) getConfigured(ctx context.Context) (settings.ThemeSettings, error) {
	set, err := s.settings.Get(ctx)
	if err != nil {
		return settings.ThemeSettings{}, fault.Wrap(err, fctx.With(ctx))
	}

	return set.Theme.OrZero(), nil
}

func (s *Service) Upload(ctx context.Context, r io.Reader, size int64, clientFilename string) (settings.ThemeAsset, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return settings.ThemeAsset{}, fault.Wrap(err, fctx.With(ctx))
	}
	if size < 1 || size > MaxAssetSizeBytes {
		return settings.ThemeAsset{}, fault.Wrap(errInvalidAsset, fmsg.Withf("theme assets must be between 1 and %d bytes", MaxAssetSizeBytes))
	}

	kind, err := classify(clientFilename)
	if err != nil {
		return settings.ThemeAsset{}, err
	}

	content, err := io.ReadAll(io.LimitReader(r, MaxAssetSizeBytes+1))
	if err != nil {
		return settings.ThemeAsset{}, fault.Wrap(err, fctx.With(ctx))
	}
	if int64(len(content)) != size || len(content) > MaxAssetSizeBytes || !utf8.Valid(content) || bytes.IndexByte(content, 0) >= 0 {
		return settings.ThemeAsset{}, fault.Wrap(errInvalidAsset, fmsg.With("theme assets must be valid UTF-8 text matching Content-Length"))
	}

	hash := sha256.Sum256(content)
	integrity := "sha256-" + base64.StdEncoding.EncodeToString(hash[:])
	filename := asset.NewFilename(clientFilename)
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return settings.ThemeAsset{}, fault.Wrap(err, fctx.With(ctx))
	}

	path := asset.BuildAssetPath(filename)
	if err := s.objects.Write(ctx, path, bytes.NewReader(content), int64(len(content))); err != nil {
		return settings.ThemeAsset{}, fault.Wrap(err, fctx.With(ctx))
	}

	stored, err := s.assets.Add(
		ctx,
		xid.ID(accountID),
		filename,
		len(content),
		mime.New(kind.MIME()),
		asset_writer.WithMetadata(
			asset.Metadata{
				metadataPurposeKey:   metadataPurposeTheme,
				metadataKindKey:      string(kind),
				metadataIntegrityKey: integrity,
			},
		),
	)
	if err != nil {
		_ = s.objects.Delete(ctx, path)
		return settings.ThemeAsset{}, fault.Wrap(err, fctx.With(ctx))
	}

	return mapAsset(stored, kind), nil
}

func (s *Service) Publish(ctx context.Context, stylesheetIDs, scriptIDs []string) (settings.ThemeSettings, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return settings.ThemeSettings{}, fault.Wrap(err, fctx.With(ctx))
	}
	if len(stylesheetIDs)+len(scriptIDs) > MaxActiveAssetCount {
		return settings.ThemeSettings{}, fault.Wrap(errInvalidTheme, fmsg.Withf("a theme may contain at most %d active assets", MaxActiveAssetCount))
	}

	stylesheets, err := s.resolve(ctx, stylesheetIDs, AssetKindStylesheet)
	if err != nil {
		return settings.ThemeSettings{}, err
	}
	scripts, err := s.resolve(ctx, scriptIDs, AssetKindScript)
	if err != nil {
		return settings.ThemeSettings{}, err
	}

	total := 0
	for _, assets := range [][]settings.ThemeAsset{stylesheets, scripts} {
		for _, item := range assets {
			total += item.Size
		}
	}
	if total > MaxActiveSizeBytes {
		return settings.ThemeSettings{}, fault.Wrap(errInvalidTheme, fmsg.Withf("active theme assets may total at most %d bytes", MaxActiveSizeBytes))
	}

	manifest := settings.ThemeSettings{
		Stylesheets: stylesheets,
		Scripts:     scripts,
	}
	_, err = s.settings.Set(ctx, settings.Settings{
		Theme: opt.New(manifest),
	})
	if err != nil {
		return settings.ThemeSettings{}, fault.Wrap(err, fctx.With(ctx))
	}

	return manifest, nil
}

func (s *Service) GetAsset(ctx context.Context, filename string) (settings.ThemeAsset, error) {
	a, err := s.assetQuery.Get(ctx, asset.NewFilepathFilename(filename))
	if err != nil {
		return settings.ThemeAsset{}, fault.Wrap(err, fctx.With(ctx))
	}
	kind, ok := themeAssetKind(a)
	if !ok {
		return settings.ThemeAsset{}, fault.Wrap(errInvalidAsset, ftag.With(ftag.NotFound))
	}
	return mapAsset(a, kind), nil
}

func (s *Service) OpenAsset(ctx context.Context, filename string) (io.Reader, int64, error) {
	r, size, err := s.objects.Read(ctx, asset.BuildAssetPath(asset.NewFilepathFilename(filename)))
	if err != nil {
		return nil, 0, fault.Wrap(err, fctx.With(ctx))
	}
	return r, size, nil
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
	for _, assets := range [][]settings.ThemeAsset{manifest.Stylesheets, manifest.Scripts} {
		for _, item := range assets {
			if item.Filename == a.Name.String() {
				return errActiveAsset
			}
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
	return nil
}

func (s *Service) resolve(ctx context.Context, ids []string, expected AssetKind) ([]settings.ThemeAsset, error) {
	seen := map[string]struct{}{}
	out := make([]settings.ThemeAsset, 0, len(ids))
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

func mapAsset(a *asset.Asset, kind AssetKind) settings.ThemeAsset {
	integrity, _ := a.Metadata[metadataIntegrityKey].(string)
	return settings.ThemeAsset{
		ID:        a.ID.String(),
		Filename:  a.Name.String(),
		MIMEType:  kind.MIME(),
		Size:      a.Size,
		Integrity: integrity,
	}
}
