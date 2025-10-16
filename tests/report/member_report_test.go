package report_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestReportCRUD(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)

			handle1 := xid.New().String()
			acc1, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{handle1, "password"})
			tests.Ok(t, err, acc1)
			session1 := sh.WithSession(e2e.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc1.JSON200.Id)))))

			handle2 := xid.New().String()
			acc2, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{handle2, "password"})
			tests.Ok(t, err, acc2)
			session2 := sh.WithSession(e2e.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc2.JSON200.Id)))))

			cat1, err := cl.CategoryCreateWithResponse(adminCtx, openapi.CategoryInitialProps{
				Colour:      "",
				Description: "cat",
				Name:        xid.New().String(),
			}, sh.WithSession(adminCtx))
			tests.Ok(t, err, cat1)

			thread1, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:     opt.New("<p>test thread content</p>").Ptr(),
				Category: &cat1.JSON200.Id,
				Title:    "test thread",
			}, session1)
			tests.Ok(t, err, thread1)

			t.Run("unauthenticated", func(t *testing.T) {
				t.Parallel()

				comment := "spam content"
				rep, err := cl.ReportCreateWithResponse(root, openapi.ReportInitialProps{
					Comment:    &comment,
					TargetId:   thread1.JSON200.Id,
					TargetKind: "post",
				})
				tests.Status(t, err, rep, http.StatusUnauthorized)
			})

			t.Run("create", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				comment := "spam content"
				rep, err := cl.ReportCreateWithResponse(root, openapi.ReportInitialProps{
					Comment:    &comment,
					TargetId:   thread1.JSON200.Id,
					TargetKind: "post",
				}, session1)
				tests.Ok(t, err, rep)
				a.Equal("spam content", *rep.JSON200.Comment)
				a.Equal(openapi.Submitted, rep.JSON200.Status)
				a.Equal(acc1.JSON200.Id, rep.JSON200.ReportedBy.Id)
				a.Nil(rep.JSON200.HandledBy)
			})

			t.Run("create_profile_report", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				comment := "inappropriate profile"
				rep, err := cl.ReportCreateWithResponse(root, openapi.ReportCreateJSONRequestBody{
					Comment:    &comment,
					TargetId:   acc2.JSON200.Id,
					TargetKind: "profile",
				}, session1)
				tests.Ok(t, err, rep)
				a.Equal("inappropriate profile", *rep.JSON200.Comment)
				a.Equal(openapi.Submitted, rep.JSON200.Status)
				a.Equal(acc1.JSON200.Id, rep.JSON200.ReportedBy.Id)
			})

			t.Run("list_own_reports", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				comment := "spam content"
				rep, err := cl.ReportCreateWithResponse(root, openapi.ReportInitialProps{
					Comment:    &comment,
					TargetId:   thread1.JSON200.Id,
					TargetKind: "post",
				}, session1)
				tests.Ok(t, err, rep)

				list1, err := cl.ReportListWithResponse(root, &openapi.ReportListParams{}, session1)
				tests.Ok(t, err, list1)

				foundReport, foundOk := lo.Find(list1.JSON200.Reports, func(r openapi.Report) bool { return r.Id == rep.JSON200.Id })
				a.True(foundOk)
				a.Equal("spam content", *foundReport.Comment)
				a.Equal(acc1.JSON200.Id, foundReport.ReportedBy.Id)

				list2, err := cl.ReportListWithResponse(root, &openapi.ReportListParams{}, session2)
				tests.Ok(t, err, list2)

				_, foundOk = lo.Find(list2.JSON200.Reports, func(r openapi.Report) bool { return r.Id == rep.JSON200.Id })
				a.False(foundOk)
			})

			t.Run("resolve_own_report", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				comment := "spam content"
				rep, err := cl.ReportCreateWithResponse(root, openapi.ReportInitialProps{
					Comment:    &comment,
					TargetId:   thread1.JSON200.Id,
					TargetKind: "post",
				}, session1)
				tests.Ok(t, err, rep)
				a.Equal(openapi.Submitted, rep.JSON200.Status)

				resolved := openapi.Resolved
				update, err := cl.ReportUpdateWithResponse(root, rep.JSON200.Id, openapi.ReportMutableProps{
					Status: &resolved,
				}, session1)
				tests.Ok(t, err, update)
				a.Equal(openapi.Resolved, update.JSON200.Status)
			})

			t.Run("cannot_resolve_others_report", func(t *testing.T) {
				t.Parallel()

				comment := "spam content"
				rep, err := cl.ReportCreateWithResponse(root, openapi.ReportInitialProps{
					Comment:    &comment,
					TargetId:   thread1.JSON200.Id,
					TargetKind: "post",
				}, session1)
				tests.Ok(t, err, rep)

				resolved := openapi.Resolved
				update, err := cl.ReportUpdateWithResponse(root, rep.JSON200.Id, openapi.ReportMutableProps{
					Status: &resolved,
				}, session2)
				tests.Status(t, err, update, http.StatusForbidden)
			})

			t.Run("cannot_change_to_acknowledged", func(t *testing.T) {
				t.Parallel()

				comment := "spam content"
				rep, err := cl.ReportCreateWithResponse(root, openapi.ReportInitialProps{
					Comment:    &comment,
					TargetId:   thread1.JSON200.Id,
					TargetKind: "post",
				}, session1)
				tests.Ok(t, err, rep)

				ack := openapi.Acknowledged
				update, err := cl.ReportUpdateWithResponse(root, rep.JSON200.Id, openapi.ReportMutableProps{
					Status: &ack,
				}, session1)
				tests.Status(t, err, update, http.StatusForbidden)
			})
		}))
	}))
}
