package report_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestReportAuthorization(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, adminAcc := e2e.WithAccount(root, aw, seed.Account_001_Odin)

			handle1 := xid.New().String()
			acc1, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{handle1, "password"})
			tests.Ok(t, err, acc1)
			session1 := sh.WithSession(e2e.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc1.JSON200.Id)))))

			handle2 := xid.New().String()
			acc2, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{handle2, "password"})
			tests.Ok(t, err, acc2)
			session2 := sh.WithSession(e2e.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc2.JSON200.Id)))))

			modRoleName := "Moderator " + xid.New().String()
			modRole, err := cl.RoleCreateWithResponse(adminCtx, openapi.RoleInitialProps{
				Name: modRoleName,
				Permissions: []openapi.Permission{
					openapi.MANAGEREPORTS,
				},
			}, sh.WithSession(adminCtx))
			tests.Ok(t, err, modRole)

			addRole1, err := cl.AccountAddRoleWithResponse(adminCtx, adminAcc.Handle, modRole.JSON200.Id, sh.WithSession(adminCtx))
			tests.Ok(t, err, addRole1)

			handle3 := xid.New().String()
			acc3, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{handle3, "password"})
			tests.Ok(t, err, acc3)
			session3 := sh.WithSession(e2e.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc3.JSON200.Id)))))

			addRole2, err := cl.AccountAddRoleWithResponse(adminCtx, handle3, modRole.JSON200.Id, sh.WithSession(adminCtx))
			tests.Ok(t, err, addRole2)

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

			t.Run("moderator_can_list_all_reports", func(t *testing.T) {
				t.Parallel()
				r := require.New(t)
				a := assert.New(t)

				comment := "spam content"
				rep, err := cl.ReportCreateWithResponse(root, openapi.ReportInitialProps{
					Comment:    &comment,
					TargetId:   thread1.JSON200.Id,
					TargetKind: "post",
				}, session1)
				tests.Ok(t, err, rep)

				findReport := func() (*openapi.Report, bool) {
					const maxPages = 20
					const maxScans = 3

					status := openapi.Submitted
					kind := "post"

					for scan := 0; scan < maxScans; scan++ {
						var page *openapi.PaginationQuery

						for p := 0; p < maxPages; p++ {
							list, err := cl.ReportListWithResponse(root, &openapi.ReportListParams{
								Page:   page,
								Status: &status,
								Kind:   &kind,
							}, session3)
							if err != nil || list == nil || list.JSON200 == nil {
								return nil, false
							}

							foundReport, foundOk := lo.Find(list.JSON200.Reports, func(r openapi.Report) bool { return r.Id == rep.JSON200.Id })
							if foundOk {
								return &foundReport, true
							}

							if list.JSON200.NextPage == nil {
								break
							}

							next := openapi.PaginationQuery(strconv.Itoa(*list.JSON200.NextPage))
							page = &next
						}
					}

					return nil, false
				}

				foundReport, foundOk := findReport()
				r.True(foundOk)
				r.NotNil(foundReport)
				a.Equal("spam content", *foundReport.Comment)
				a.Equal(acc1.JSON200.Id, foundReport.ReportedBy.Id)
			})

			t.Run("moderator_can_acknowledge_report", func(t *testing.T) {
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

				ack := openapi.Acknowledged
				update, err := cl.ReportUpdateWithResponse(root, rep.JSON200.Id, openapi.ReportMutableProps{
					Status: &ack,
				}, session3)
				tests.Ok(t, err, update)
				a.Equal(openapi.Acknowledged, update.JSON200.Status)
				a.NotNil(update.JSON200.HandledBy)
				a.Equal(acc3.JSON200.Id, update.JSON200.HandledBy.Id)
			})

			t.Run("moderator_can_resolve_report", func(t *testing.T) {
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
				}, session3)
				tests.Ok(t, err, update)
				a.Equal(openapi.Resolved, update.JSON200.Status)
				a.NotNil(update.JSON200.HandledBy)
				a.Equal(acc3.JSON200.Id, update.JSON200.HandledBy.Id)
			})

			t.Run("moderator_can_assign_handler", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				comment := "spam content"
				rep, err := cl.ReportCreateWithResponse(root, openapi.ReportInitialProps{
					Comment:    &comment,
					TargetId:   thread1.JSON200.Id,
					TargetKind: "post",
				}, session1)
				tests.Ok(t, err, rep)
				a.Nil(rep.JSON200.HandledBy)

				ack := openapi.Acknowledged
				adminID := adminAcc.ID.String()
				update, err := cl.ReportUpdateWithResponse(root, rep.JSON200.Id, openapi.ReportMutableProps{
					Status:    &ack,
					HandledBy: &adminID,
				}, session3)
				tests.Ok(t, err, update)
				a.Equal(openapi.Acknowledged, update.JSON200.Status)
				a.NotNil(update.JSON200.HandledBy)
				a.Equal(adminID, update.JSON200.HandledBy.Id)
			})

			t.Run("regular_member_cannot_list_all_reports", func(t *testing.T) {
				t.Parallel()
				r := require.New(t)
				a := assert.New(t)

				comment1 := "spam 1"
				rep1, err := cl.ReportCreateWithResponse(root, openapi.ReportInitialProps{
					Comment:    &comment1,
					TargetId:   thread1.JSON200.Id,
					TargetKind: "post",
				}, session1)
				tests.Ok(t, err, rep1)

				comment2 := "spam 2"
				rep2, err := cl.ReportCreateWithResponse(root, openapi.ReportInitialProps{
					Comment:    &comment2,
					TargetId:   thread1.JSON200.Id,
					TargetKind: "post",
				}, session2)
				tests.Ok(t, err, rep2)

				list1, err := cl.ReportListWithResponse(root, &openapi.ReportListParams{}, session1)
				tests.Ok(t, err, list1)

				_, foundRep1 := lo.Find(list1.JSON200.Reports, func(r openapi.Report) bool { return r.Id == rep1.JSON200.Id })
				r.True(foundRep1)

				_, foundRep2 := lo.Find(list1.JSON200.Reports, func(r openapi.Report) bool { return r.Id == rep2.JSON200.Id })
				a.False(foundRep2)
			})

			t.Run("regular_member_cannot_acknowledge_report", func(t *testing.T) {
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

			t.Run("filter_by_status", func(t *testing.T) {
				t.Parallel()
				r := require.New(t)
				a := assert.New(t)

				comment1 := "spam 1"
				rep1, err := cl.ReportCreateWithResponse(root, openapi.ReportInitialProps{
					Comment:    &comment1,
					TargetId:   thread1.JSON200.Id,
					TargetKind: "post",
				}, session1)
				tests.Ok(t, err, rep1)

				comment2 := "spam 2"
				rep2, err := cl.ReportCreateWithResponse(root, openapi.ReportInitialProps{
					Comment:    &comment2,
					TargetId:   thread1.JSON200.Id,
					TargetKind: "post",
				}, session1)
				tests.Ok(t, err, rep2)

				ack := openapi.Acknowledged
				update2, err := cl.ReportUpdateWithResponse(root, rep2.JSON200.Id, openapi.ReportMutableProps{
					Status: &ack,
				}, session3)
				tests.Ok(t, err, update2)

				submittedStatus := openapi.Submitted
				listSubmitted, err := cl.ReportListWithResponse(root, &openapi.ReportListParams{
					Status: &submittedStatus,
				}, session3)
				tests.Ok(t, err, listSubmitted)

				_, foundRep1 := lo.Find(listSubmitted.JSON200.Reports, func(r openapi.Report) bool { return r.Id == rep1.JSON200.Id })
				r.True(foundRep1)

				_, foundRep2 := lo.Find(listSubmitted.JSON200.Reports, func(r openapi.Report) bool { return r.Id == rep2.JSON200.Id })
				a.False(foundRep2)

				acknowledgedStatus := openapi.Acknowledged
				listAcknowledged, err := cl.ReportListWithResponse(root, &openapi.ReportListParams{
					Status: &acknowledgedStatus,
				}, session3)
				tests.Ok(t, err, listAcknowledged)

				_, foundRep1Ack := lo.Find(listAcknowledged.JSON200.Reports, func(r openapi.Report) bool { return r.Id == rep1.JSON200.Id })
				a.False(foundRep1Ack)

				_, foundRep2Ack := lo.Find(listAcknowledged.JSON200.Reports, func(r openapi.Report) bool { return r.Id == rep2.JSON200.Id })
				r.True(foundRep2Ack)
			})

			t.Run("filter_by_kind", func(t *testing.T) {
				t.Parallel()
				r := require.New(t)
				a := assert.New(t)

				commentPost := "spam post"
				repPost, err := cl.ReportCreateWithResponse(root, openapi.ReportInitialProps{
					Comment:    &commentPost,
					TargetId:   thread1.JSON200.Id,
					TargetKind: "post",
				}, session1)
				tests.Ok(t, err, repPost)

				commentProfile := "bad profile"
				repProfile, err := cl.ReportCreateWithResponse(root, openapi.ReportCreateJSONRequestBody{
					Comment:    &commentProfile,
					TargetId:   acc2.JSON200.Id,
					TargetKind: "profile",
				}, session1)
				tests.Ok(t, err, repProfile)

				kindPost := "post"
				listPost, err := cl.ReportListWithResponse(root, &openapi.ReportListParams{
					Kind: &kindPost,
				}, session3)
				tests.Ok(t, err, listPost)

				_, foundPost := lo.Find(listPost.JSON200.Reports, func(r openapi.Report) bool { return r.Id == repPost.JSON200.Id })
				r.True(foundPost)

				_, foundProfile := lo.Find(listPost.JSON200.Reports, func(r openapi.Report) bool { return r.Id == repProfile.JSON200.Id })
				a.False(foundProfile)

				kindProfile := "profile"
				listProfile, err := cl.ReportListWithResponse(root, &openapi.ReportListParams{
					Kind: &kindProfile,
				}, session3)
				tests.Ok(t, err, listProfile)

				_, foundPostInProfile := lo.Find(listProfile.JSON200.Reports, func(r openapi.Report) bool { return r.Id == repPost.JSON200.Id })
				a.False(foundPostInProfile)

				_, foundProfileInProfile := lo.Find(listProfile.JSON200.Reports, func(r openapi.Report) bool { return r.Id == repProfile.JSON200.Id })
				r.True(foundProfileInProfile)
			})
		}))
	}))
}
