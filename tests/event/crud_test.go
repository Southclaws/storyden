package event_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/event/location"
	"github.com/Southclaws/storyden/app/resources/event/participation"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestEventsCRUD(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *session_cookie.Jar,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			adminCtx, adminAcc := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := e2e.WithSession(adminCtx, cj)

			cat, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Admin: false, Colour: "#fe4efd", Description: "category testing", Name: "Category " + uuid.NewString()}, adminSession)
			tests.Ok(t, err, cat)
			catID := cat.JSON200.Id

			t.Run("create", func(t *testing.T) {
				name := "New Cool Event " + uuid.NewString()
				timerange := openapi.EventTimeRange{Start: time.Now().Add(time.Hour * 72), End: time.Now().Add(time.Hour * 96)}
				create, err := cl.EventCreateWithResponse(root, openapi.EventInitialProps{
					Name:                name,
					Description:         opt.New("description of my lovely event").Ptr(),
					Content:             "<body><p>hello world</p></body>",
					TimeRange:           timerange,
					ParticipationPolicy: openapi.Open,
					Visibility:          openapi.Published,
					Capacity:            opt.New(14).Ptr(),
					ThreadCategoryId:    catID,
				}, adminSession)
				tests.Ok(t, err, create)

				a.NotEmpty(create.JSON200.Id)
				a.NotEmpty(create.JSON200.CreatedAt)
				a.NotEmpty(create.JSON200.UpdatedAt)
				a.Empty(create.JSON200.DeletedAt)

				a.Equal(name, create.JSON200.Name)
				a.Contains(create.JSON200.Slug, "new-cool-event")
				a.Equal("description of my lovely event", create.JSON200.Description)
				a.WithinDuration(create.JSON200.TimeRange.Start, timerange.Start, time.Second*5)
				a.WithinDuration(create.JSON200.TimeRange.End, timerange.End, time.Second*5)
				a.Equal(openapi.Open, create.JSON200.ParticipationPolicy)
				a.Equal(openapi.Published, create.JSON200.Visibility)
				matchLocation(t, &location.Virtual{}, create.JSON200.Location)
				a.Equal(14, *create.JSON200.Capacity)
				a.Equal("<body><p>hello world</p></body>", create.JSON200.Thread.Body)

				r.Len(create.JSON200.Participants, 1)
				a.Equal(adminAcc.ID.String(), create.JSON200.Participants[0].Profile.Id)
				a.Equal(participation.RoleHost.String(), string(create.JSON200.Participants[0].Role))
			})
		}))
	}))
}

func matchLocation(t *testing.T, want location.Location, got openapi.EventLocation) {
	t.Helper()
	r := require.New(t)
	a := assert.New(t)

	switch w := want.(type) {
	case *location.Physical:
		g, err := got.AsEventLocationPhysical()
		r.NoError(err)
		r.NotNil(g)

		a.Equal(w.Name, g.Name)
		a.Equal(w.Address, opt.NewPtr(g.Address))
		a.Equal(w.Latitude, opt.NewPtr(g.Latitude))
		a.Equal(w.Longitude, opt.NewPtr(g.Longitude))
		a.Equal(w.URL, opt.NewPtr(g.Url))

	case *location.Virtual:
		g, err := got.AsEventLocationVirtual()
		r.NoError(err)
		r.NotNil(g)

		a.Equal(w.Name, g.Name)
		a.Equal(w.URL.String(), opt.NewPtr(g.Url).String())
	}
}
