package asset_migration

import (
	"time"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
)

const Kind = "asset_migration"

type Target string

const (
	Avatar Target = "avatar"
	Icon   Target = "icon"
	Banner Target = "banner"
)

type Payload struct {
	Target           Target            `json:"target"`
	Path             string            `json:"path"`
	AccountID        account.AccountID `json:"account_id"`
	AccountUpdatedAt time.Time         `json:"account_updated_at,omitempty"`
	SourceID         asset.AssetID     `json:"source_id,omitempty"`
}

func (p Payload) Key() string {
	if p.Target == Avatar {
		return p.Path
	}

	return string(p.Target)
}

func (p Payload) Presets() int64 {
	presets := asset.EagerImagePresets

	switch p.Target {
	case Icon:
		presets |= asset.IconImagePresets
	case Banner:
		presets |= asset.BannerImagePresets
	}

	return presets
}

func LegacyFilenames(target Target) []string {
	switch target {
	case Icon:
		return []string{
			"icon-512x512.png",
			"icon-180x180.png",
			"icon-167x167.png",
			"icon-152x152.png",
			"icon-120x120.png",
			"icon-32x32.png",
		}
	case Banner:
		return []string{"banner-1200x630.png"}
	default:
		return nil
	}
}
