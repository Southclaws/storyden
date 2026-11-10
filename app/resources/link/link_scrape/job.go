package link_scrape

import (
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
)

const Kind = "link_scrape"

type Payload struct {
	URL  link_url.URL   `json:"url"`
	Item *datagraph.Ref `json:"item,omitempty"`
}
