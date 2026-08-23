package bindings

import (
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/trail"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func TestSerialiseTrailRunWithoutTriggerDetails(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.August, 23, 18, 0, 0, 0, time.UTC)
	run := &trail.Run{
		ID:        trail.RunID(xid.New()),
		TrailID:   trail.ID(xid.New()),
		Kind:      trail.RunKindEvent,
		Status:    trail.RunStatusCompleted,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := serialiseTrailRun(run)
	require.NoError(t, err)
	require.Nil(t, result.Trigger)
	require.Equal(t, openapi.TrailRunStatusCompleted, result.Status)
}
