package recurrence

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreview(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		raw  string
		now  string
		want []string
	}{
		"one time through count": {
			raw:  `{"start":"2026-09-03T09:30:00","timezone":"Europe/London","rule":{"frequency":"daily","interval":1,"count":1}}`,
			now:  "2026-09-01T00:00:00Z",
			want: []string{"2026-09-03T08:30:00Z"},
		},
		"every two weeks": {
			raw:  `{"start":"2026-08-18T10:00:00","timezone":"Europe/London","rule":{"frequency":"weekly","interval":2,"by_weekday":["tuesday"]}}`,
			now:  "2026-08-18T09:00:00Z",
			want: []string{"2026-09-01T09:00:00Z", "2026-09-15T09:00:00Z"},
		},
		"start anchors selectors without matching them": {
			raw:  `{"start":"2026-08-23T09:00:00","timezone":"Europe/London","rule":{"frequency":"weekly","interval":1,"by_weekday":["tuesday","thursday"]}}`,
			now:  "2026-08-22T00:00:00Z",
			want: []string{"2026-08-25T08:00:00Z", "2026-08-27T08:00:00Z", "2026-09-01T08:00:00Z"},
		},
		"yearly Christmas Eve": {
			raw:  `{"start":"2026-12-24T18:00:00","timezone":"Europe/London","rule":{"frequency":"yearly","interval":1,"by_month":[12],"by_month_day":[24]}}`,
			now:  "2026-01-01T00:00:00Z",
			want: []string{"2026-12-24T18:00:00Z", "2027-12-24T18:00:00Z"},
		},
		"monthly last day": {
			raw:  `{"start":"2026-01-31T12:15:00","timezone":"UTC","rule":{"frequency":"monthly","interval":1,"by_month_day":[-1]}}`,
			now:  "2026-01-30T13:00:00Z",
			want: []string{"2026-01-31T12:15:00Z", "2026-02-28T12:15:00Z", "2026-03-31T12:15:00Z"},
		},
		"invalid month days are skipped": {
			raw:  `{"start":"2026-01-31T12:15:00","timezone":"UTC","rule":{"frequency":"monthly","interval":1,"by_month_day":[31]}}`,
			now:  "2026-01-31T12:15:00Z",
			want: []string{"2026-03-31T12:15:00Z", "2026-05-31T12:15:00Z"},
		},
		"daily keeps local time over DST": {
			raw:  `{"start":"2026-03-27T09:00:00","timezone":"Europe/London","rule":{"frequency":"daily","interval":1}}`,
			now:  "2026-03-28T09:01:00Z",
			want: []string{"2026-03-29T08:00:00Z", "2026-03-30T08:00:00Z"},
		},
		"hourly skips a DST gap": {
			raw:  `{"start":"2026-03-29T00:30:00","timezone":"Europe/London","rule":{"frequency":"hourly","interval":1}}`,
			now:  "2026-03-29T00:30:00Z",
			want: []string{"2026-03-29T01:30:00Z", "2026-03-29T02:30:00Z"},
		},
		"yearly skips invalid leap days": {
			raw:  `{"start":"2024-02-29T12:00:00","timezone":"UTC","rule":{"frequency":"yearly","interval":1,"by_month":[2],"by_month_day":[29]}}`,
			now:  "2024-02-29T12:00:00Z",
			want: []string{"2028-02-29T12:00:00Z", "2032-02-29T12:00:00Z"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			schedule, err := Parse(json.RawMessage(tc.raw))
			require.NoError(t, err)
			now, err := time.Parse(time.RFC3339, tc.now)
			require.NoError(t, err)

			got := Preview(schedule, now, len(tc.want))
			formatted := make([]string, len(got))
			for i, occurrence := range got {
				formatted[i] = occurrence.Format(time.RFC3339)
			}
			assert.Equal(t, tc.want, formatted)
		})
	}
}

func TestSkipsNonexistentLocalTime(t *testing.T) {
	t.Parallel()

	schedule, err := Parse(json.RawMessage(`{"start":"2026-03-22T01:30:00","timezone":"Europe/London","rule":{"frequency":"weekly","interval":1,"by_weekday":["sunday"]}}`))
	require.NoError(t, err)
	now, err := time.Parse(time.RFC3339, "2026-03-22T02:00:00Z")
	require.NoError(t, err)

	got := Preview(schedule, now, 2)
	require.Len(t, got, 2)
	assert.Equal(t, "2026-04-05T00:30:00Z", got[0].Format(time.RFC3339))
	assert.Equal(t, "2026-04-12T00:30:00Z", got[1].Format(time.RFC3339))
}

func TestFiresAmbiguousTimeOnce(t *testing.T) {
	t.Parallel()

	schedule, err := Parse(json.RawMessage(`{"start":"2026-10-18T01:30:00","timezone":"Europe/London","rule":{"frequency":"weekly","interval":1,"by_weekday":["sunday"]}}`))
	require.NoError(t, err)
	now, err := time.Parse(time.RFC3339, "2026-10-24T00:00:00Z")
	require.NoError(t, err)

	got := Preview(schedule, now, 2)
	require.Len(t, got, 2)
	assert.Equal(t, "2026-10-25T01:30:00Z", got[0].Format(time.RFC3339))
	assert.Equal(t, "2026-11-01T01:30:00Z", got[1].Format(time.RFC3339))
}

func TestNormalisesSelectorOrder(t *testing.T) {
	t.Parallel()

	schedule, err := Parse(json.RawMessage(`{"start":"2026-12-24T18:00:00","timezone":"UTC","rule":{"frequency":"yearly","interval":1,"by_month":[12,1],"by_month_day":[24,1]}}`))
	require.NoError(t, err)
	encoded, err := schedule.NormalisedJSON()
	require.NoError(t, err)
	assert.JSONEq(t, `{"start":"2026-12-24T18:00:00","timezone":"UTC","rule":{"frequency":"yearly","interval":1,"by_month":[1,12],"by_month_day":[1,24]}}`, string(encoded))
}

func TestValidate(t *testing.T) {
	t.Parallel()

	invalid := []string{
		`{"start":"2026-09-03T09:30:00","timezone":"Not/AZone","rule":{"frequency":"daily","interval":1}}`,
		`{"start":"2026-09-03T09:30:00Z","timezone":"UTC","rule":{"frequency":"daily","interval":1}}`,
		`{"start":"2026-09-03T09:30:00","timezone":"UTC","rule":{"frequency":"daily","interval":0}}`,
		`{"start":"2026-09-03T09:30:00","timezone":"UTC","rule":{"frequency":"weekly","interval":1,"by_weekday":["thursday","thursday"]}}`,
		`{"start":"2026-09-03T09:30:00","timezone":"UTC","rule":{"frequency":"weekly","interval":1,"by_month_day":[3]}}`,
		`{"start":"2026-09-03T09:30:00","timezone":"UTC","rule":{"frequency":"monthly","interval":1,"by_month_day":[0]}}`,
		`{"start":"2026-09-03T09:30:00","timezone":"UTC","rule":{"frequency":"yearly","interval":1,"by_month_day":[3]}}`,
		`{"start":"2026-09-03T09:30:00","timezone":"UTC","rule":{"frequency":"daily","interval":1},"extra":true}`,
	}
	for _, raw := range invalid {
		_, err := Parse(json.RawMessage(raw))
		assert.Error(t, err, raw)
	}
}
