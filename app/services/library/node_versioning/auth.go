package node_versioning

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/library/node_version"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func authoriseDraftMutation(ctx context.Context, callerID account.AccountID, v *node_version.NodeVersion) error {
	if v.Status != node_version.VersionStatusDraft {
		return fault.New("version is not a draft",
			fctx.With(ctx),
			ftag.With(ftag.PermissionDenied),
			fmsg.WithDesc("not a draft", "Only draft versions can be modified."),
		)
	}

	return session.Authorise(ctx, func() error {
		if v.Author.ID != callerID {
			return fault.New("not the draft author",
				fctx.With(ctx),
				ftag.With(ftag.PermissionDenied),
				fmsg.WithDesc("not author", "You can only modify your own drafts."),
			)
		}

		return nil
	}, rbac.PermissionManageLibrary)
}

func authoriseDraftVisible(ctx context.Context, callerID account.AccountID, v *node_version.NodeVersion) error {
	err := session.Authorise(ctx, func() error {
		if v.Author.ID != callerID {
			return fault.New("not the draft author",
				fctx.With(ctx),
				ftag.With(ftag.PermissionDenied),
				fmsg.WithDesc("not author", "You can only view your own draft versions."),
			)
		}

		return nil
	}, rbac.PermissionManageLibrary)
	if err != nil && ftag.Get(err) == ftag.PermissionDenied {
		return versionNotFound(ctx)
	}

	return err
}

func authoriseVisibleDraftMutation(ctx context.Context, callerID account.AccountID, v *node_version.NodeVersion) error {
	err := authoriseDraftMutation(ctx, callerID, v)
	if err != nil && v.Status == node_version.VersionStatusDraft && ftag.Get(err) == ftag.PermissionDenied {
		return versionNotFound(ctx)
	}

	return err
}

func authoriseDraftDelete(ctx context.Context, callerID account.AccountID, v *node_version.NodeVersion) error {
	if v.Status != node_version.VersionStatusDraft {
		return fault.New("version is not a draft",
			fctx.With(ctx),
			ftag.With(ftag.PermissionDenied),
			fmsg.WithDesc("not a draft", "Only draft versions can be discarded."),
		)
	}

	return session.Authorise(ctx, func() error {
		if v.Author.ID != callerID {
			return fault.New("not the draft author",
				fctx.With(ctx),
				ftag.With(ftag.PermissionDenied),
				fmsg.WithDesc("not author", "You can only discard your own drafts."),
			)
		}

		return nil
	}, rbac.PermissionManageLibrary)
}

func authoriseVisibleDraftDelete(ctx context.Context, callerID account.AccountID, v *node_version.NodeVersion) error {
	err := authoriseDraftDelete(ctx, callerID, v)
	if err != nil && v.Status == node_version.VersionStatusDraft && ftag.Get(err) == ftag.PermissionDenied {
		return versionNotFound(ctx)
	}

	return err
}
