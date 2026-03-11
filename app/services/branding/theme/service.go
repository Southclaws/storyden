package theme

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_writer"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/asset/asset_download"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/mime"
)

var (
	errInvalidThemeAssetName = fault.Wrap(
		fault.New("invalid theme asset filename"),
		ftag.With(ftag.InvalidArgument),
	)
	errInvalidThemeManifest = fault.Wrap(
		fault.New("invalid theme manifest"),
		ftag.With(ftag.InvalidArgument),
	)
	errInvalidThemeAssetType = fault.Wrap(
		fault.New("invalid theme asset mime type"),
		ftag.With(ftag.InvalidArgument),
	)
)

const (
	ThemeMetadataKey = "theme"
)

type AssetKind string

const (
	AssetKindStylesheet AssetKind = "css"
	AssetKindScript     AssetKind = "js"
)

func (k AssetKind) MIME() string {
	switch k {
	case AssetKindStylesheet:
		return "text/css"
	case AssetKindScript:
		return "application/javascript"
	default:
		return "application/octet-stream"
	}
}

func (k AssetKind) Extension() string {
	switch k {
	case AssetKindStylesheet:
		return ".css"
	case AssetKindScript:
		return ".js"
	default:
		return ""
	}
}

type Manifest struct {
	CSS     []string `json:"css"`
	Scripts []string `json:"scripts"`
}

func NewManifest(css, scripts []string) Manifest {
	return Manifest{
		CSS:     css,
		Scripts: scripts,
	}
}

func DefaultManifest() Manifest {
	return Manifest{
		CSS:     []string{},
		Scripts: []string{},
	}
}

type Service interface {
	GetManifest(ctx context.Context) (Manifest, error)
	UploadAsset(ctx context.Context, r io.Reader, size int64, clientFilename string) (*asset.Asset, AssetKind, error)
	GetAsset(ctx context.Context, filename string) (*asset.Asset, io.Reader, error)
}

type service struct {
	settings   *settings.SettingsRepository
	assets     *asset_writer.Writer
	downloader *asset_download.Downloader
	objects    object.Storer
}

func New(
	settings *settings.SettingsRepository,
	assets *asset_writer.Writer,
	downloader *asset_download.Downloader,
	objects object.Storer,
) Service {
	return &service{
		settings:   settings,
		assets:     assets,
		downloader: downloader,
		objects:    objects,
	}
}

func (s *service) GetManifest(ctx context.Context) (Manifest, error) {
	set, err := s.settings.Get(ctx)
	if err != nil {
		return DefaultManifest(), fault.Wrap(err, fctx.With(ctx))
	}

	manifest, err := ParseManifest(set.Metadata)
	if err != nil {
		return DefaultManifest(), fault.Wrap(err, fctx.With(ctx))
	}

	return manifest, nil
}

func (s *service) UploadAsset(
	ctx context.Context,
	r io.Reader,
	size int64,
	clientFilename string,
) (*asset.Asset, AssetKind, error) {
	kind, err := classifyAssetKind(clientFilename)
	if err != nil {
		return nil, "", fault.Wrap(err, fctx.With(ctx))
	}

	detectedMIME, recycled, err := mime.Detect(r)
	if err != nil {
		return nil, "", fault.Wrap(err, fctx.With(ctx))
	}
	if err := validateDetectedMIME(kind, detectedMIME.String()); err != nil {
		return nil, "", fault.Wrap(err, fctx.With(ctx))
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, "", fault.Wrap(err, fctx.With(ctx))
	}

	filename := asset.NewFilename(clientFilename)
	assetMIME := mime.New(kind.Extension())

	a, err := s.assets.Add(ctx, xid.ID(accountID), filename, int(size), assetMIME)
	if err != nil {
		return nil, "", fault.Wrap(err, fctx.With(ctx))
	}

	path := asset.BuildAssetPath(a.Name)
	if err := s.objects.Write(ctx, path, recycled, size); err != nil {
		return nil, "", fault.Wrap(err, fctx.With(ctx))
	}

	return a, kind, nil
}

func (s *service) GetAsset(ctx context.Context, filename string) (*asset.Asset, io.Reader, error) {
	a, r, err := s.downloader.Get(ctx, asset.NewFilepathFilename(filename))
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	return a, r, nil
}

func ParseManifest(in opt.Optional[map[string]any]) (Manifest, error) {
	metadata := in.Or(map[string]any{})

	rawTheme, exists := metadata[ThemeMetadataKey]
	if !exists {
		return DefaultManifest(), nil
	}

	themeMap, ok := rawTheme.(map[string]any)
	if !ok {
		return DefaultManifest(), fault.Wrap(
			errInvalidThemeManifest,
			fmsg.With("theme metadata must be an object"),
		)
	}

	css := parseStringList(themeMap["css"])
	scripts := parseStringList(themeMap["scripts"])

	css = dedupeStrings(css)
	scripts = dedupeStrings(scripts)

	return NewManifest(css, scripts), nil
}

func BuildAssetURL(filename string) string {
	return fmt.Sprintf("/api/info/theme/assets/%s", asset.NewFilepathFilename(filename).String())
}

func normalizeManifestAssetPath(v string) string {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "/") {
		return trimmed
	}

	parsed, err := url.Parse(trimmed)
	if err == nil && parsed.IsAbs() {
		return trimmed
	}

	return BuildAssetURL(trimmed)
}

func parseStringList(v any) []string {
	if v == nil {
		return []string{}
	}

	switch raw := v.(type) {
	case []string:
		out := make([]string, 0, len(raw))
		for _, item := range raw {
			path := normalizeManifestAssetPath(item)
			if path != "" {
				out = append(out, path)
			}
		}
		return out

	case []any:
		out := make([]string, 0, len(raw))
		for _, item := range raw {
			s, ok := item.(string)
			if !ok {
				continue
			}
			path := normalizeManifestAssetPath(s)
			if path != "" {
				out = append(out, path)
			}
		}
		return out

	default:
		return []string{}
	}
}

func dedupeStrings(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	out := make([]string, 0, len(input))

	for _, item := range input {
		if _, exists := seen[item]; exists {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}

	return out
}

func classifyAssetKind(clientFilename string) (AssetKind, error) {
	trimmed := strings.TrimSpace(clientFilename)
	if trimmed == "" {
		return "", fault.Wrap(
			errInvalidThemeAssetName,
			fmsg.With("filename is required"),
		)
	}

	if strings.Contains(trimmed, "/") || strings.Contains(trimmed, `\`) {
		return "", fault.Wrap(
			errInvalidThemeAssetName,
			fmsg.With("filename must not include path separators"),
		)
	}

	ext := strings.ToLower(filepath.Ext(trimmed))
	switch ext {
	case ".css":
		return AssetKindStylesheet, nil
	case ".js":
		return AssetKindScript, nil
	default:
		return "", fault.Wrap(
			errInvalidThemeAssetName,
			fmsg.With("theme assets must use .css or .js filename extensions"),
		)
	}
}

func validateDetectedMIME(kind AssetKind, detected string) error {
	base := strings.TrimSpace(strings.SplitN(strings.ToLower(detected), ";", 2)[0])

	switch kind {
	case AssetKindStylesheet:
		if base == "text/css" || base == "text/plain" {
			return nil
		}
	case AssetKindScript:
		if base == "application/javascript" || base == "text/javascript" || base == "text/plain" {
			return nil
		}
	}

	return fault.Wrap(
		errInvalidThemeAssetType,
		fmsg.Withf("detected mime type %q is not valid for %s theme assets", detected, kind),
	)
}
