package writer

import (
	"context"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/rs/xid"
)

type Writer struct {
	client    *ent.Client
	batchSize int

	RoleIDMap     map[int]xid.ID
	AccountIDMap  map[int]xid.ID
	CategoryIDMap map[int]xid.ID
	PostIDMap     map[int]xid.ID
	TagIDMap      map[int]xid.ID
	FirstPostPIDs map[int]bool
}

func New(client *ent.Client, batchSize int) *Writer {
	return &Writer{
		client:        client,
		batchSize:     batchSize,
		RoleIDMap:     make(map[int]xid.ID),
		AccountIDMap:  make(map[int]xid.ID),
		CategoryIDMap: make(map[int]xid.ID),
		PostIDMap:     make(map[int]xid.ID),
		TagIDMap:      make(map[int]xid.ID),
	}
}

func (w *Writer) Client() *ent.Client {
	return w.client
}

func (w *Writer) BatchSize() int {
	return w.batchSize
}

func (w *Writer) CreateRoles(ctx context.Context, builders []*ent.RoleCreate) ([]*ent.Role, error) {
	return w.client.Role.CreateBulk(builders...).Save(ctx)
}

func (w *Writer) CreateAccounts(ctx context.Context, builders []*ent.AccountCreate) ([]*ent.Account, error) {
	var result []*ent.Account
	for i := 0; i < len(builders); i += w.batchSize {
		end := i + w.batchSize
		if end > len(builders) {
			end = len(builders)
		}
		batch, err := w.client.Account.CreateBulk(builders[i:end]...).Save(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, batch...)
	}
	return result, nil
}

func (w *Writer) CreateAuthentications(ctx context.Context, builders []*ent.AuthenticationCreate) ([]*ent.Authentication, error) {
	var result []*ent.Authentication
	for i := 0; i < len(builders); i += w.batchSize {
		end := i + w.batchSize
		if end > len(builders) {
			end = len(builders)
		}
		batch, err := w.client.Authentication.CreateBulk(builders[i:end]...).Save(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, batch...)
	}
	return result, nil
}

func (w *Writer) CreateAccountRoles(ctx context.Context, builders []*ent.AccountRolesCreate) ([]*ent.AccountRoles, error) {
	var result []*ent.AccountRoles
	for i := 0; i < len(builders); i += w.batchSize {
		end := i + w.batchSize
		if end > len(builders) {
			end = len(builders)
		}
		batch, err := w.client.AccountRoles.CreateBulk(builders[i:end]...).Save(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, batch...)
	}
	return result, nil
}

func (w *Writer) CreateEmails(ctx context.Context, builders []*ent.EmailCreate) ([]*ent.Email, error) {
	var result []*ent.Email
	for i := 0; i < len(builders); i += w.batchSize {
		end := i + w.batchSize
		if end > len(builders) {
			end = len(builders)
		}
		batch, err := w.client.Email.CreateBulk(builders[i:end]...).Save(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, batch...)
	}
	return result, nil
}

func (w *Writer) CreateCategories(ctx context.Context, builders []*ent.CategoryCreate) ([]*ent.Category, error) {
	return w.client.Category.CreateBulk(builders...).Save(ctx)
}

func (w *Writer) CreateTags(ctx context.Context, builders []*ent.TagCreate) ([]*ent.Tag, error) {
	return w.client.Tag.CreateBulk(builders...).Save(ctx)
}

func (w *Writer) CreatePosts(ctx context.Context, builders []*ent.PostCreate) ([]*ent.Post, error) {
	var result []*ent.Post
	for i := 0; i < len(builders); i += w.batchSize {
		end := i + w.batchSize
		if end > len(builders) {
			end = len(builders)
		}
		batch, err := w.client.Post.CreateBulk(builders[i:end]...).Save(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, batch...)
	}
	return result, nil
}

func (w *Writer) CreateReacts(ctx context.Context, builders []*ent.ReactCreate) ([]*ent.React, error) {
	return w.client.React.CreateBulk(builders...).Save(ctx)
}

func (w *Writer) CreateLikePosts(ctx context.Context, builders []*ent.LikePostCreate) ([]*ent.LikePost, error) {
	return w.client.LikePost.CreateBulk(builders...).Save(ctx)
}

func (w *Writer) CreatePostReads(ctx context.Context, builders []*ent.PostReadCreate) ([]*ent.PostRead, error) {
	return w.client.PostRead.CreateBulk(builders...).Save(ctx)
}

func (w *Writer) CreateReports(ctx context.Context, builders []*ent.ReportCreate) ([]*ent.Report, error) {
	return w.client.Report.CreateBulk(builders...).Save(ctx)
}

func (w *Writer) CreateAssets(ctx context.Context, builders []*ent.AssetCreate) ([]*ent.Asset, error) {
	return w.client.Asset.CreateBulk(builders...).Save(ctx)
}
