package asset_image

import "github.com/Southclaws/storyden/app/resources/asset"

type Payload struct {
	Requested int64                    `json:"requested"`
	Completed int64                    `json:"completed"`
	Skipped   map[asset.Variant]string `json:"skipped,omitempty"`
}

func (p Payload) IsComplete() bool {
	return p.Requested == p.Completed
}
