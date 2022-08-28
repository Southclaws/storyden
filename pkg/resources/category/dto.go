package category

import (
	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	"github.com/rs/xid"
)

type CategoryID xid.ID

type PostMeta struct {
	Author string    `json:"author"`
	PostID xid.ID `json:"postId"`
	Slug   string    `json:"slug"`
	Title  string    `json:"title"`
	Short  string    `json:"short"`
}

type Category struct {
	ID          CategoryID `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Colour      string     `json:"colour"`
	Sort        int        `json:"sort"`
	Admin       bool       `json:"admin"`
	Recent      []PostMeta `json:"recent,omitempty"`
	PostCount   int        `json:"postCount"`
}

func PostMetaFromModel(p *model.Post) *PostMeta {
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

func FromModel(c *model.Category) *Category {
	recent := []PostMeta{}
	for _, p := range c.Edges.Posts {
		recent = append(recent, *PostMetaFromModel(p))
	}

	return &Category{
		ID:          CategoryID(c.ID),
		Name:        c.Name,
		Description: c.Description,
		Colour:      c.Colour,
		Sort:        c.Sort,
		Admin:       c.Admin,
		Recent:      recent,
	}
}
