package link_check

import (
	"context"
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
	ent_link "github.com/Southclaws/storyden/internal/ent/link"
)

const (
	Kind             = "link_liveness"
	OfflineThreshold = 3
	MinimumAge       = 90 * 24 * time.Hour
	CheckInterval    = 365 * 24 * time.Hour
	RetryInterval    = 30 * 24 * time.Hour
)

type Status string

const (
	Unknown Status = "unknown"
	Alive   Status = "alive"
	Offline Status = "offline"
)

type Candidate struct {
	ID       xid.ID
	URL      string
	Failures int
}

type Repository struct {
	db *ent.Client
}

func New(db *ent.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	return r.db.Link.Query().Count(ctx)
}

// Due reads raw stored URLs so legacy policy failures can be removed.
func (r *Repository) Due(ctx context.Context, now time.Time, limit int) ([]Candidate, error) {
	if limit <= 0 {
		return nil, nil
	}

	query := r.db.Link.Query().Where(ent_link.CreatedAtLTE(now.Add(-MinimumAge))).
		Select(ent_link.FieldID, ent_link.FieldURL, ent_link.FieldLivenessFailures, ent_link.FieldLivenessCheckedAt).
		Order(ent_link.ByLivenessCheckedAt(), ent_link.ByID())

	unchecked, err := query.Clone().Where(ent_link.LivenessCheckedAtIsNil()).Limit(limit).All(ctx)
	if err != nil {
		return nil, err
	}

	checked, err := query.Where(
		ent_link.LivenessCheckedAtLTE(now.Add(-RetryInterval)),
		ent_link.Or(
			ent_link.LivenessCheckedAtLTE(now.Add(-CheckInterval)),
			ent_link.And(
				ent_link.LivenessFailuresGT(0),
				ent_link.LivenessFailuresLT(OfflineThreshold),
			),
		),
	).Limit(limit).All(ctx)
	if err != nil {
		return nil, err
	}

	uncheckedQuota := limit / 2
	// Alternate the extra slot so one-link batches give both groups a turn.
	if now.UTC().YearDay()%2 == 1 {
		uncheckedQuota = (limit + 1) / 2
	}

	uncheckedCount := min(uncheckedQuota, len(unchecked))
	checkedCount := min(limit-uncheckedCount, len(checked))
	uncheckedCount = min(limit-checkedCount, len(unchecked))
	rows := append(unchecked[:uncheckedCount], checked[:checkedCount]...)

	out := make([]Candidate, 0, len(rows))
	for _, row := range rows {
		out = append(out, Candidate{ID: row.ID, URL: row.URL, Failures: row.LivenessFailures})
	}

	return out, nil
}

func (r *Repository) Record(ctx context.Context, candidate Candidate, status Status, now time.Time) error {
	failures := 0
	if status == Offline {
		failures = candidate.Failures + 1
		if failures < OfflineThreshold {
			status = Unknown
		}
	}

	err := r.db.Link.UpdateOneID(candidate.ID).
		SetLivenessStatus(ent_link.LivenessStatus(status)).
		SetLivenessCheckedAt(now).
		SetLivenessFailures(failures).
		Exec(ctx)
	if ent.IsNotFound(err) {
		return nil
	}

	return err
}

func (r *Repository) Remove(ctx context.Context, id xid.ID) error {
	err := r.db.Link.DeleteOneID(id).Exec(ctx)
	if ent.IsNotFound(err) {
		return nil
	}

	return err
}
