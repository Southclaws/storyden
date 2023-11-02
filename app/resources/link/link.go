package link

import "github.com/Southclaws/opt"

type Link struct {
	URL         string
	Title       opt.Optional[string]
	Description opt.Optional[string]
}

func NewLink(url, title, description string) Link {
	return Link{
		URL:         url,
		Title:       opt.New(title),
		Description: opt.New(description),
	}
}

func NewLinkOpt(purl, ptitle, pdescription *string) opt.Optional[Link] {
	if purl == nil {
		return opt.NewEmpty[Link]()
	}

	return opt.New(Link{
		URL:         opt.NewPtr(purl).String(),
		Title:       opt.NewPtr(ptitle),
		Description: opt.NewPtr(pdescription),
	})
}
