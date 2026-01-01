package writer

import (
	"context"
	"fmt"

	"github.com/Southclaws/storyden/cmd/import-mybb/logger"
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

func (w *Writer) DeleteAllData(ctx context.Context) (err error) {
	var i int

	i, err = w.client.Post.Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete Posts: %w", err)
	}
	logger.Success(fmt.Sprintf("Deleted %d existing Posts (Threads and Replies)", i))

	i, err = w.client.Account.Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete Accounts: %w", err)
	}
	logger.Success(fmt.Sprintf("Deleted %d existing Accounts", i))

	i, err = w.client.Role.Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete Roles: %w", err)
	}
	logger.Success(fmt.Sprintf("Deleted %d existing Roles", i))

	i, err = w.client.Authentication.Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete Authentications: %w", err)
	}
	logger.Success(fmt.Sprintf("Deleted %d existing Authentications", i))

	i, err = w.client.AccountRoles.Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete AccountRoless: %w", err)
	}
	logger.Success(fmt.Sprintf("Deleted %d existing AccountRoless", i))

	i, err = w.client.Email.Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete Emails: %w", err)
	}
	logger.Success(fmt.Sprintf("Deleted %d existing Emails", i))

	i, err = w.client.Category.Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete Categorys: %w", err)
	}
	logger.Success(fmt.Sprintf("Deleted %d existing Categories", i))

	i, err = w.client.Tag.Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete Tags: %w", err)
	}
	logger.Success(fmt.Sprintf("Deleted %d existing Tags", i))

	i, err = w.client.React.Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete Reacts: %w", err)
	}
	logger.Success(fmt.Sprintf("Deleted %d existing Reacts", i))

	i, err = w.client.LikePost.Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete LikePosts: %w", err)
	}
	logger.Success(fmt.Sprintf("Deleted %d existing LikePosts", i))

	i, err = w.client.PostRead.Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete PostReads: %w", err)
	}
	logger.Success(fmt.Sprintf("Deleted %d existing PostReads", i))

	i, err = w.client.Report.Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete Reports: %w", err)
	}
	logger.Success(fmt.Sprintf("Deleted %d existing Reports", i))

	i, err = w.client.Asset.Delete().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete Assets: %w", err)
	}
	logger.Success(fmt.Sprintf("Deleted %d existing Assets", i))

	return nil
}

func (w *Writer) CreateRoles(ctx context.Context, builders []*ent.RoleCreate) ([]*ent.Role, error) {
	r, err := w.client.Role.CreateBulk(builders...).Save(ctx)
	if err != nil {
		return nil, err
	}

	logger.Success(fmt.Sprintf("Imported %d Roles", len(r)))

	return r, nil
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

	logger.Success(fmt.Sprintf("Imported %d Accounts", len(result)))

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

	logger.Success(fmt.Sprintf("Imported %d Authentications", len(result)))

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

	logger.Success(fmt.Sprintf("Imported %d AccountRoless", len(result)))

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

	logger.Success(fmt.Sprintf("Imported %d Emails", len(result)))

	return result, nil
}

func (w *Writer) CreateCategories(ctx context.Context, builders []*ent.CategoryCreate) ([]*ent.Category, error) {
	r, err := w.client.Category.CreateBulk(builders...).Save(ctx)
	if err != nil {
		return nil, err
	}

	logger.Success(fmt.Sprintf("Imported %d Categories", len(r)))

	return r, nil
}

func (w *Writer) CreateTags(ctx context.Context, builders []*ent.TagCreate) ([]*ent.Tag, error) {
	r, err := w.client.Tag.CreateBulk(builders...).Save(ctx)
	if err != nil {
		return nil, err
	}

	logger.Success(fmt.Sprintf("Imported %d Tags", len(r)))

	return r, nil
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

	logger.Success(fmt.Sprintf("Imported %d Posts (Threads and Replies)", len(result)))

	return result, nil
}

func (w *Writer) CreateReacts(ctx context.Context, builders []*ent.ReactCreate) ([]*ent.React, error) {
	r, err := w.client.React.CreateBulk(builders...).Save(ctx)
	if err != nil {
		return nil, err
	}

	logger.Success(fmt.Sprintf("Imported %d React", len(r)))

	return r, nil
}

func (w *Writer) CreateLikePosts(ctx context.Context, builders []*ent.LikePostCreate) ([]*ent.LikePost, error) {
	r, err := w.client.LikePost.CreateBulk(builders...).Save(ctx)
	if err != nil {
		return nil, err
	}

	logger.Success(fmt.Sprintf("Imported %d LikePost", len(r)))

	return r, nil
}

func (w *Writer) CreatePostReads(ctx context.Context, builders []*ent.PostReadCreate) ([]*ent.PostRead, error) {
	r, err := w.client.PostRead.CreateBulk(builders...).Save(ctx)
	if err != nil {
		return nil, err
	}

	logger.Success(fmt.Sprintf("Imported %d PostRead", len(r)))

	return r, nil
}

func (w *Writer) CreateReports(ctx context.Context, builders []*ent.ReportCreate) ([]*ent.Report, error) {
	r, err := w.client.Report.CreateBulk(builders...).Save(ctx)
	if err != nil {
		return nil, err
	}

	logger.Success(fmt.Sprintf("Imported %d Reports", len(r)))

	return r, nil
}

func (w *Writer) CreateAssets(ctx context.Context, builders []*ent.AssetCreate) ([]*ent.Asset, error) {
	r, err := w.client.Asset.CreateBulk(builders...).Save(ctx)
	if err != nil {
		return nil, err
	}

	logger.Success(fmt.Sprintf("Imported %d Assets", len(r)))

	return r, nil
}
