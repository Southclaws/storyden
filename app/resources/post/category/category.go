package category

import (
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/internal/ent"
)

type CategoryID xid.ID

func (i CategoryID) String() string { return xid.ID(i).String() }

type PostMeta struct {
	Author string
	PostID xid.ID
	Slug   string
	Title  string
	Short  string
}

type Category struct {
	ID          CategoryID
	Name        string
	Slug        string
	Description string
	Colour      string
	Sort        int
	Admin       bool
	ParentID    *CategoryID
	CoverImage  opt.Optional[asset.Asset]
	Children    []*Category
	Recent      []PostMeta
	PostCount   int
	Metadata    map[string]any
	UpdatedAt   time.Time
}

func PostMetaFromModel(p *ent.Post) *PostMeta {
	slug := p.Slug

	title := p.Title

	return &PostMeta{
		Author: p.Edges.Author.Name,
		PostID: p.ID,
		Slug:   slug,
		Title:  title,
		Short:  p.Short,
	}
}

func FromModel(c *ent.Category) *Category {
	recent := []PostMeta{}

	if c.Edges.Posts != nil {
		for _, p := range c.Edges.Posts {
			recent = append(recent, *PostMetaFromModel(p))
		}
	}

	var parentID *CategoryID

	if !c.ParentCategoryID.IsNil() {
		pid := CategoryID(c.ParentCategoryID)
		parentID = &pid
	}

	coverImage := opt.Map(opt.NewPtr(c.Edges.CoverImage), func(a ent.Asset) asset.Asset {
		return *asset.Map(&a)
	})

	children := dt.Map(c.Edges.Children, FromModel)

	return &Category{
		ID:          CategoryID(c.ID),
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		Colour:      c.Colour,
		Sort:        c.Sort,
		Admin:       c.Admin,
		ParentID:    parentID,
		CoverImage:  coverImage,
		Children:    children,
		Recent:      recent,
		Metadata:    c.Metadata,
		UpdatedAt:   c.UpdatedAt,
	}
}
