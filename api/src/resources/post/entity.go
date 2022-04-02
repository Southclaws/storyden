package post

import (
	"fmt"
	"time"

	"github.com/Southclaws/storyden/api/src/infra/db"
	"github.com/Southclaws/storyden/api/src/resources/category"
	"github.com/Southclaws/storyden/api/src/resources/react"
)

type Author struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Admin     bool      `json:"admin"`
	CreatedAt time.Time `json:"createdAt"`
}

type Post struct {
	ID         string             `json:"id"`
	Title      *string            `json:"title"`
	Slug       *string            `json:"slug"`
	Body       string             `json:"body"`
	Short      string             `json:"short"`
	First      bool               `json:"first"`
	Pinned     bool               `json:"pinned"`
	CreatedAt  time.Time          `json:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt"`
	DeletedAt  *time.Time         `json:"deletedAt"`
	UserID     string             `json:"userId"`
	RootPostID *string            `json:"rootPostId"`
	Author     Author             `json:"author"`
	Tags       []string           `json:"tags"`
	Posts      int                `json:"posts"`
	ReplyTo    *category.PostMeta `json:"replyTo"`
	Category   category.Category  `json:"category"`
	Reacts     []react.React      `json:"reacts"`
}

func FromModel(u *db.PostModel) (w *Post) {
	tags := []string{}
	if u.RelationsPost.Tags != nil {
		for _, t := range u.RelationsPost.Tags {
			tags = append(tags, t.Name)
		}
	}

	cat := category.Category{}
	if u.RelationsPost.Category != nil {
		cat = *category.FromModel(u.RelationsPost.Category)
	}

	var replyTo *category.PostMeta
	if u.RelationsPost.ReplyTo != nil {
		replyTo = &category.PostMeta{
			Author: u.RelationsPost.Author.Name,
			PostID: u.RelationsPost.ReplyTo.ID,
		}
	}

	reacts := []react.React{}
	if u.RelationsPost.Reacts != nil {
		for _, r := range u.RelationsPost.Reacts {
			reacts = append(reacts, *react.FromModel(&r, u.ID))
		}
	}

	return &Post{
		ID:         u.InnerPost.ID,
		Slug:       getSlug(u),
		Title:      u.InnerPost.Title,
		Body:       u.InnerPost.Body,
		Short:      u.InnerPost.Short,
		First:      u.InnerPost.First,
		Pinned:     u.InnerPost.Pinned,
		CreatedAt:  u.InnerPost.CreatedAt,
		UpdatedAt:  u.InnerPost.UpdatedAt,
		DeletedAt:  u.InnerPost.DeletedAt,
		UserID:     u.InnerPost.UserID,
		RootPostID: u.InnerPost.RootPostID,
		Author: Author{
			ID:        u.RelationsPost.Author.InnerUser.ID,
			Name:      u.RelationsPost.Author.InnerUser.Name,
			Admin:     u.RelationsPost.Author.InnerUser.Admin,
			CreatedAt: u.RelationsPost.Author.InnerUser.CreatedAt,
		},
		Tags:     tags,
		Category: cat,
		ReplyTo:  replyTo,
		Reacts:   reacts,
	}
}

func getSlug(post *db.PostModel) *string {
	slug, ok := post.Slug()
	if ok {
		return &slug
	}

	if post.RelationsPost.Root != nil && post.RelationsPost.Root.InnerPost.Slug != nil {
		slug = fmt.Sprintf("%s#%s", *post.RelationsPost.Root.InnerPost.Slug, post.ID)
	}

	return &slug
}
