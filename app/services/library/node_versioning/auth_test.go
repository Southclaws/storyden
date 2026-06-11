package node_versioning

import (
	"context"
	"testing"

	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/library/node_version"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func TestAuthoriseDraftMutation(t *testing.T) {
	t.Parallel()

	callerID := account.AccountID(xid.New())
	otherID := account.AccountID(xid.New())

	for _, tc := range []struct {
		name    string
		version *node_version.NodeVersion
		wantErr bool
		manager bool
	}{
		{
			name: "non-draft denied",
			version: &node_version.NodeVersion{
				Status: node_version.VersionStatusApplied,
				Author: profile.Ref{ID: callerID},
			},
			wantErr: true,
		},
		{
			name: "different author denied",
			version: &node_version.NodeVersion{
				Status: node_version.VersionStatusDraft,
				Author: profile.Ref{ID: otherID},
			},
			wantErr: true,
		},
		{
			name: "draft author allowed",
			version: &node_version.NodeVersion{
				Status: node_version.VersionStatusDraft,
				Author: profile.Ref{ID: callerID},
			},
			wantErr: false,
		},
		{
			name: "manager allowed",
			version: &node_version.NodeVersion{
				Status: node_version.VersionStatusDraft,
				Author: profile.Ref{ID: otherID},
			},
			manager: true,
			wantErr: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			a := assert.New(t)

			permissions := rbac.NewList(rbac.PermissionSubmitLibraryNodeChanges)
			if tc.manager {
				permissions = rbac.NewList(rbac.PermissionManageLibrary)
			}

			ctx := session.WithAccountPermissions(context.Background(), account.Account{ID: callerID}, permissions)

			err := authoriseDraftMutation(ctx, callerID, tc.version)

			if tc.wantErr {
				a.Error(err)
				a.Equal(ftag.PermissionDenied, ftag.Get(err))
				return
			}

			a.NoError(err)
		})
	}
}
