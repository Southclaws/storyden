package thread

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"

	"github.com/Southclaws/storyden/api/src/infra/db"
	"github.com/Southclaws/storyden/api/src/infra/db/model"
	"github.com/Southclaws/storyden/api/src/resources/post"
	"github.com/Southclaws/storyden/api/src/resources/user"
)

var (
	ErrNoTitle      = errors.New("missing title")
	ErrNoBody       = errors.New("missing body")
	ErrUnauthorised = errors.New("unauthorised")
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) CreateThread(
	ctx context.Context,
	title, body string,
	authorID user.UserID,
	categoryName string,
	tags []string,
) (*post.Post, error) {
	slug := slug.Make(title)
	short := post.MakeShortBody(body)

	tagset, err := d.createTags(ctx, tags)
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert tags for linking to post")
	}

	d.db.Post.
		Create().
		SetBody(body).
		SetShort(short).
		SetFirst(true).
		SetAuthorID(uuid.UUID(authorID)).
		SetTitle(title).
		// SetCategoryID(categoryName).
		AddTagIDs(tagset).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	p, err = d.db.Post.
		FindUnique(
			db.Post.ID.Equals(p.ID),
		).
		With(
			db.Post.Author.Fetch(),
			db.Post.Category.Fetch(),
			db.Post.Tags.Fetch(),
		).
		Update(
			db.Post.Slug.Set(fmt.Sprintf("%s-%s", p.ID, slug)),
		).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return post.FromModel(p), nil
}

func (d *database) createTags(ctx context.Context, tags []string) ([]db.TagWhereParam, error) {
	setters := []db.TagWhereParam{}
	for _, tag := range tags {
		if len(tag) > 24 {
			return nil, post.ErrTagNameTooLong
		}
		_, err := d.db.Tag.
			UpsertOne(db.Tag.Name.Equals(tag)).
			Update().
			Create(db.Tag.Name.Set(tag)).
			Exec(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to upsert tag")
		}
		setters = append(setters, db.Tag.Name.Equals(tag))
	}
	return setters, nil
}

func (d *database) CreateLegacyThread(
	ctx context.Context,
	title string,
	categoryName string,
	tags []string,
	first post.Post,
	replies []post.Post,
) (*post.Post, error) {
	slug := slug.Make(title)
	short := post.MakeShortBody(first.Body)

	var firstUpdatedAt *time.Time
	if !first.UpdatedAt.IsZero() {
		firstUpdatedAt = &first.UpdatedAt
	} else {
		firstUpdatedAt = &first.CreatedAt
	}

	tagset, err := d.createTags(ctx, tags)
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert tags for linking to post")
	}

	additional := []db.PostSetParam{
		db.Post.Title.Set(title),
		db.Post.Category.Link(db.Category.Name.Equals(categoryName)),
		db.Post.Tags.Link(tagset...),

		db.Post.CreatedAt.Set(first.CreatedAt),
		db.Post.UpdatedAt.SetIfPresent(firstUpdatedAt),
	}

	p, err := d.db.Post.
		CreateOne(
			db.Post.Body.Set(first.Body),
			db.Post.Short.Set(short),
			db.Post.First.Set(true),
			db.Post.Author.Link(db.User.ID.Equals(first.Author.ID)),

			additional...,
		).
		With(db.Post.Author.Fetch()).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	_, err = d.db.Post.
		FindUnique(db.Post.ID.Equals(p.ID)).
		Update(
			db.Post.Slug.Set(fmt.Sprintf("%s-%s", p.ID, slug)),
			db.Post.UpdatedAt.SetIfPresent(&first.CreatedAt), // keep the updated date old
		).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	parentID := p.ID

	for _, r := range replies {
		short := post.MakeShortBody(r.Body)

		var updatedAt *time.Time
		if !r.UpdatedAt.IsZero() {
			updatedAt = &r.UpdatedAt
		} else {
			updatedAt = &r.CreatedAt
		}

		optional := []db.PostSetParam{
			db.Post.Root.Link(db.Post.ID.Equals(parentID)),

			db.Post.CreatedAt.Set(r.CreatedAt),
			db.Post.UpdatedAt.SetIfPresent(updatedAt),
		}

		_, err = d.db.Post.
			CreateOne(
				db.Post.Body.Set(r.Body),
				db.Post.Short.Set(short),
				db.Post.First.Set(false),
				db.Post.Author.Link(db.User.ID.Equals(r.Author.ID)),

				optional...,
			).
			With(db.Post.Author.Fetch()).
			Exec(ctx)
		if err != nil {
			return nil, err
		}

	}

	return post.FromModel(p), nil
}

func (d *database) GetThreads(
	ctx context.Context,
	tags []string,
	category string,
	query string,
	before time.Time,
	sort string,
	offset int,
	max int,
	includePosts bool,
	includeDeleted bool,
	includeAdmin bool,
) ([]post.Post, error) {
	filters := []db.PostWhereParam{}

	if !before.IsZero() {
		filters = append(filters, db.Post.CreatedAt.Before(before))
	}
	if category != "" {
		filters = append(filters, db.Post.Category.Where(db.Category.Name.Equals(category)))
	}
	if query != "" {
		filters = append(filters, db.Post.Title.Search(query), db.Post.Body.Search(query))
	}
	if len(tags) > 0 {
		for _, v := range tags {
			if len(v) > 0 {
				filters = append(filters, db.Post.Tags.Some(db.Tag.Name.Equals(v)))
			}
		}
	}
	if max < 1 {
		max = 1
	}
	if max > 100 {
		max = 100
	}
	if !includePosts {
		filters = append(filters, db.Post.First.Equals(true))
	}
	if !includeDeleted {
		filters = append(filters, db.Post.DeletedAt.IsNull())
	}
	if !includeAdmin {
		filters = append(filters, db.Post.Category.Where(
			db.Category.Admin.Equals(false),
		))
	}

	var pinned []db.PostModel
	var err error
	if category != "" {
		pinned, err = d.db.Post.
			FindMany(append(filters, db.Post.Pinned.Equals(true))...).
			With(
				db.Post.Author.Fetch(),
				db.Post.Tags.Fetch(),
				db.Post.Category.Fetch(),
			).
			Skip(offset).
			Take(max).
			OrderBy(db.Post.CreatedAt.Order(db.Direction(sort))).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	posts, err := d.db.Post.
		FindMany(filters...).
		With(
			db.Post.Author.Fetch(),
			db.Post.Tags.Fetch(),
			db.Post.Category.Fetch(),
		).
		Skip(offset).
		Take(max).
		OrderBy(db.Post.CreatedAt.Order(db.Direction(sort))).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	counts, err := d.GetPostCounts(ctx)
	if err != nil {
		return nil, err
	}

	//
	// NOTE:
	//
	// This double query is a temporary solution for pinned posts. The reason
	// for this is the multiple sort order query didn't work, details here:
	// https://github.com/prisma/prisma-client-go/issues/681
	//
	// NOTE 2:
	//
	// This code was removed, but the removal was reverted because of this bug:
	// https://github.com/prisma/prisma-client-go/issues/698
	//

	result := []post.Post{}
	pinnedIDs := make(map[string]struct{})
	// First, add the pinned posts to the start of the list.
	for _, p := range pinned {
		p.Body = ""
		newpost := *post.FromModel(&p)
		newpost.Posts = counts[p.ID]
		result = append(result, newpost)

		// store the pinned post ID so we can filter it out below.
		pinnedIDs[p.ID] = struct{}{}
	}
	// then add the rest of the posts, filtering out what's already there.
	for _, p := range posts {
		// don't duplicate pinned posts they are already at the top of the list.
		if _, ok := pinnedIDs[p.ID]; ok {
			continue
		}
		p.Body = ""
		newpost := *post.FromModel(&p)
		newpost.Posts = counts[p.ID]
		result = append(result, newpost)
	}

	return result, nil
}

func (d *database) GetPostCounts(ctx context.Context) (map[string]int, error) {
	type PostCount struct {
		PostID string `json:"rootPostId"`
		Count  int    `json:"count"`
	}

	var counts []PostCount
	err := d.db.Prisma.Raw.QueryRaw(`
		with recursive counts AS(
			select id, "rootPostId"
			from public."Post"
			where "rootPostId" is not null

			union

			select s.id, s."rootPostId"
			from public."Post" s
			inner join counts c on c.id = s."rootPostId"
		) select "rootPostId", count(*) from counts
		group by "rootPostId"`).
		Exec(ctx, &counts)
	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, c := range counts {
		result[c.PostID] = c.Count
	}

	return result, nil
}

func (d *database) Update(ctx context.Context, userID, id string, title, categoryID *string, pinned *bool) (*post.Post, error) {
	updates := []db.PostSetParam{
		db.Post.Title.SetIfPresent(title),
		db.Post.CategoryID.SetIfPresent(categoryID),
		db.Post.Pinned.SetIfPresent(pinned),
	}

	if err := post.CanUserMutatePost(ctx, d.db, userID, id); err != nil {
		return nil, post.ErrUnauthorised
	}

	p, err := d.db.Post.
		FindUnique(db.Post.ID.Equals(id)).
		With(
			db.Post.Author.Fetch(),
			db.Post.Category.Fetch(),
			db.Post.Tags.Fetch(),
		).
		Update(updates...).
		Exec(ctx)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return post.FromModel(p), nil
}

func (d *database) Delete(ctx context.Context, id, authorID string) (int, error) {
	// NOTE:
	// We really want this authorID to eventually be removed from this API.
	// Because this API should be user-agnostic, and should be usable by non-
	// human users. Therefore, the validation of access rights should happen at
	// a different abstraction layer. Lower than the HTTP API but higher than
	// the database implementation.
	if err := post.CanUserMutatePost(ctx, d.db, authorID, id); err != nil {
		return 0, errors.Wrap(err, "failed to check user permissions")
	}

	result, err := d.db.Post.FindMany(
		db.Post.Or(
			db.Post.And(
				db.Post.First.Equals(true),
				db.Post.ID.Equals(id),
			),
			db.Post.And(
				db.Post.First.Equals(false),
				db.Post.Root.Where(db.Post.ID.Equals(id)),
			),
		),
	).Update(
		db.Post.DeletedAt.Set(time.Now()),
	).Exec(ctx)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return 0, nil
		}
		return 0, errors.Wrap(err, "failed to set deletedAt for posts")
	}

	return result.Count, nil
}
