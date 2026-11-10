package link_ref

import (
	"sort"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/Southclaws/storyden/internal/ent"
)

type ID = xid.ID

type LinkRef struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time

	URL               link_url.URL
	Slug              string
	Domain            string
	Title             opt.Optional[string]
	Description       opt.Optional[string]
	LivenessStatus    string
	LivenessCheckedAt *time.Time
	ScrapeStatus      string
	ScrapeAttemptedAt *time.Time
	ScrapedAt         *time.Time
	FaviconImage      opt.Optional[asset.Asset]
	PrimaryImage      opt.Optional[asset.Asset]
}

type LinkRefs []*LinkRef

func (l LinkRefs) Latest() opt.Optional[*LinkRef] {
	if len(l) == 0 {
		return opt.NewEmpty[*LinkRef]()
	}

	return opt.New(l[0])
}

func (a LinkRefs) Len() int           { return len(a) }
func (a LinkRefs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LinkRefs) Less(i, j int) bool { return xid.ID(a[i].ID).String() > xid.ID(a[j].ID).String() }

// LinksFromModel omits legacy URLs that no longer satisfy the Link policy.
func LinksFromModel(in []*ent.Link) []*LinkRef {
	list := make([]*LinkRef, 0, len(in))
	for _, row := range in {
		ref, err := Map(row)
		if err != nil {
			continue
		}

		list = append(list, ref)
	}

	sort.Sort(LinkRefs(list))
	return list
}

func Map(in *ent.Link) (*LinkRef, error) {
	u, err := link_url.Parse(in.URL)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	updatedAt := in.CreatedAt
	if in.ScrapedAt != nil && in.ScrapedAt.After(updatedAt) {
		updatedAt = *in.ScrapedAt
	}

	favicon := opt.NewPtrMap(in.Edges.FaviconImage, func(a ent.Asset) asset.Asset {
		return *asset.Map(&a)
	})

	primary := opt.NewPtrMap(in.Edges.PrimaryImage, func(a ent.Asset) asset.Asset {
		return *asset.Map(&a)
	})

	return &LinkRef{
		ID:                ID(in.ID),
		CreatedAt:         in.CreatedAt,
		UpdatedAt:         updatedAt,
		URL:               u,
		Slug:              in.Slug,
		Domain:            in.Domain,
		Title:             opt.New(in.Title),
		Description:       opt.New(in.Description),
		LivenessStatus:    string(in.LivenessStatus),
		LivenessCheckedAt: in.LivenessCheckedAt,
		ScrapeStatus:      string(in.ScrapeStatus),
		ScrapeAttemptedAt: in.ScrapeAttemptedAt,
		ScrapedAt:         in.ScrapedAt,
		FaviconImage:      favicon,
		PrimaryImage:      primary,
	}, nil
}

// MapOptional omits missing links and legacy URLs rejected by the current policy.
func MapOptional(in *ent.Link) opt.Optional[LinkRef] {
	if in == nil {
		return opt.NewEmpty[LinkRef]()
	}

	ref, err := Map(in)
	if err != nil {
		return opt.NewEmpty[LinkRef]()
	}

	return opt.New(*ref)
}
