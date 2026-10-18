package background_job

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/ent"
	ent_backgroundjob "github.com/Southclaws/storyden/internal/ent/backgroundjob"
)

func TestMapJobDecodesTypedPayload(t *testing.T) {
	t.Parallel()

	type payload struct {
		Sequence int64 `json:"sequence"`
	}

	id := xid.New()
	expires := time.Date(2026, time.October, 17, 12, 0, 0, 0, time.UTC)
	job, err := Map[payload](&ent.BackgroundJob{
		ID:         id,
		Kind:       "test.sequence",
		Status:     ent_backgroundjob.StatusProcessing,
		LeaseToken: "worker token",
		LeaseUntil: expires,
		Payload:    json.RawMessage(`{"sequence":1152921504606846977}`),
	})
	require.NoError(t, err)
	require.Equal(t, id, job.ID)
	require.Equal(t, StatusProcessing, job.Status)
	require.Equal(t, "worker token", job.LeaseToken)
	require.Equal(t, expires, job.LeaseUntil)
	require.Equal(t, int64(1152921504606846977), job.Payload.Sequence)
}

func TestMapJobRejectsInvalidRecords(t *testing.T) {
	t.Parallel()

	type payload struct {
		Sequence int64 `json:"sequence"`
	}

	for _, test := range []struct {
		name   string
		status ent_backgroundjob.Status
		raw    json.RawMessage
	}{
		{name: "unknown status", status: "unknown", raw: json.RawMessage(`{"sequence":1}`)},
		{name: "wrong payload type", status: ent_backgroundjob.StatusPending, raw: json.RawMessage(`{"sequence":"invalid"}`)},
	} {
		t.Run(test.name, func(t *testing.T) {
			job, err := Map[payload](&ent.BackgroundJob{Status: test.status, Payload: test.raw})
			require.Error(t, err)
			require.Nil(t, job)
		})
	}
}
