package action_dispatcher

import (
	"context"
	"errors"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/audit"
	"github.com/Southclaws/storyden/app/resources/audit/audit_writer"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
	ent_collection "github.com/Southclaws/storyden/internal/ent/collection"
	ent_like "github.com/Southclaws/storyden/internal/ent/likepost"
	ent_node "github.com/Southclaws/storyden/internal/ent/node"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	ent_react "github.com/Southclaws/storyden/internal/ent/react"
)

func Build() fx.Option {
	return fx.Provide(New)
}

type Service struct {
	db            *ent.Client
	accountWriter *account_writer.Writer
	auditWriter   *audit_writer.Writer
}

func New(
	db *ent.Client,
	accountWriter *account_writer.Writer,
	auditWriter *audit_writer.Writer,
) *Service {
	return &Service{
		db:            db,
		accountWriter: accountWriter,
		auditWriter:   auditWriter,
	}
}

func (s *Service) PurgeAccountContent(
	ctx context.Context,
	accountID account.AccountID,
	enactedBy opt.Optional[account.AccountID],
	contentTypes []ContentType,
) (*audit.AuditLog, error) {
	if len(contentTypes) == 0 {
		return nil, fault.New("no content types specified for account content purge", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	var errs error

	for _, ct := range contentTypes {
		switch ct {
		case ContentTypeThreads:
			if err := s.purgeThreads(ctx, accountID); err != nil {
				errs = errors.Join(errs, fault.Wrap(err, fctx.With(ctx)))
			}

		case ContentTypeReplies:
			if err := s.purgeReplies(ctx, accountID); err != nil {
				errs = errors.Join(errs, fault.Wrap(err, fctx.With(ctx)))
			}

		case ContentTypeReacts:
			if err := s.purgeReacts(ctx, accountID); err != nil {
				errs = errors.Join(errs, fault.Wrap(err, fctx.With(ctx)))
			}

		case ContentTypeLikes:
			if err := s.purgeLikes(ctx, accountID); err != nil {
				errs = errors.Join(errs, fault.Wrap(err, fctx.With(ctx)))
			}

		case ContentTypeNodes:
			if err := s.purgeNodes(ctx, accountID); err != nil {
				errs = errors.Join(errs, fault.Wrap(err, fctx.With(ctx)))
			}

		case ContentTypeCollections:
			if err := s.purgeCollections(ctx, accountID); err != nil {
				errs = errors.Join(errs, fault.Wrap(err, fctx.With(ctx)))
			}

		case ContentTypeProfileBio:
			if err := s.purgeProfileBio(ctx, accountID); err != nil {
				errs = errors.Join(errs, fault.Wrap(err, fctx.With(ctx)))
			}
		}
	}

	contentTypeStrings := dt.Map(contentTypes, func(ct ContentType) string {
		return ct.String()
	})

	ref := opt.New(datagraph.Ref{
		ID:   xid.ID(accountID),
		Kind: datagraph.KindProfile,
	})
	meta := map[string]any{
		"account_id": accountID.String(),
		"included":   contentTypeStrings,
	}

	var auditLog *audit.AuditLog
	var err error
	if errs != nil {
		auditLog, err = s.auditWriter.RecordFailure(ctx, audit.EventTypeAccountContentPurged, enactedBy, ref, meta, errs)
	} else {
		auditLog, err = s.auditWriter.Create(ctx, audit.EventTypeAccountContentPurged, enactedBy, ref, meta)
	}
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return auditLog, nil
}

func (s *Service) purgeThreads(ctx context.Context, accountID account.AccountID) error {
	// Bulk soft-delete threads using Ent update
	_, err := s.db.Post.Update().
		Where(
			ent_post.AccountPosts(xid.ID(accountID)),
			ent_post.DeletedAtIsNil(),
			ent_post.RootPostIDIsNil(),
		).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return nil
}

func (s *Service) purgeReplies(ctx context.Context, accountID account.AccountID) error {
	// Bulk soft-delete replies using Ent update
	_, err := s.db.Post.Update().
		Where(
			ent_post.AccountPosts(xid.ID(accountID)),
			ent_post.DeletedAtIsNil(),
			ent_post.RootPostIDNotNil(),
		).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return nil
}

func (s *Service) purgeReacts(ctx context.Context, accountID account.AccountID) error {
	_, err := s.db.React.Delete().
		Where(ent_react.AccountID(xid.ID(accountID))).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return nil
}

func (s *Service) purgeLikes(ctx context.Context, accountID account.AccountID) error {
	_, err := s.db.LikePost.Delete().
		Where(ent_like.AccountID(xid.ID(accountID))).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return nil
}

func (s *Service) purgeNodes(ctx context.Context, accountID account.AccountID) error {
	// Bulk soft-delete nodes using Ent update
	_, err := s.db.Node.Update().
		Where(
			ent_node.AccountID(xid.ID(accountID)),
			ent_node.DeletedAtIsNil(),
		).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return nil
}

func (s *Service) purgeCollections(ctx context.Context, accountID account.AccountID) error {
	_, err := s.db.Collection.Delete().
		Where(ent_collection.HasOwnerWith(ent_account.ID(xid.ID(accountID)))).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return nil
}

func (s *Service) purgeProfileBio(ctx context.Context, accountID account.AccountID) error {
	_, err := s.accountWriter.Update(ctx, accountID, account_writer.SetBio(""))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
