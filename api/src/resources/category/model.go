package category

import "github.com/Southclaws/storyden/api/src/infra/db"

type PostMeta struct {
	Author string `json:"author"`
	PostID string `json:"postId"`
	Slug   string `json:"slug"`
	Title  string `json:"title"`
	Short  string `json:"short"`
}

type Category struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Colour      string     `json:"colour"`
	Sort        int        `json:"sort"`
	Admin       bool       `json:"admin"`
	Recent      []PostMeta `json:"recent,omitempty"`
	PostCount   int        `json:"postCount"`
}

func PostMetaFromModel(p *db.PostModel) *PostMeta {
	slug, ok := p.Slug()
	if !ok {
		slug = ""
	}
	title, ok := p.Title()
	if !ok {
		title = ""
	}
	return &PostMeta{
		Author: p.RelationsPost.Author.Name,
		PostID: p.ID,
		Slug:   slug,
		Title:  title,
		Short:  p.Short,
	}
}

func FromModel(c *db.CategoryModel) *Category {
	recent := []PostMeta{}
	for _, p := range c.RelationsCategory.Posts {
		recent = append(recent, *PostMetaFromModel(&p))
	}

	return &Category{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Colour:      c.Colour,
		Sort:        c.Sort,
		Admin:       c.Admin,
		Recent:      recent,
	}
}
