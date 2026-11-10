package liveness

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	_ "github.com/glebarez/go-sqlite"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	jobs "github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/app/resources/link/link_check"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/Southclaws/storyden/app/resources/link/link_writer"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/enttest"
	entlink "github.com/Southclaws/storyden/internal/ent/link"
)

type probeFunc func(context.Context, link_url.URL) (link_check.Status, error)

func (f probeFunc) Check(ctx context.Context, u link_url.URL) (link_check.Status, error) {
	return f(ctx, u)
}

func newTestChecker(t *testing.T) (*ent.Client, *Checker) {
	t.Helper()
	raw, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "links.db")+"?_pragma=foreign_keys(1)")
	require.NoError(t, err)
	raw.SetMaxOpenConns(1)
	client := enttest.NewClient(t, enttest.WithOptions(ent.Driver(entsql.OpenDB(dialect.SQLite, raw))))
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	return client, &Checker{
		links:  link_check.New(client),
		jobs:   jobs.New(client),
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		prober: probeFunc(func(context.Context, link_url.URL) (link_check.Status, error) {
			return link_check.Alive, nil
		}),
	}
}

func storedLink(t *testing.T, db *ent.Client, raw string) *ent.Link {
	t.Helper()
	old := time.Now().UTC().Add(-2 * link_check.CheckInterval).Truncate(time.Microsecond)
	row, err := db.Link.Create().SetURL(raw).SetSlug(xid.New().String()).SetDomain("example.com").
		SetTitle("Saved title").SetDescription("Saved description").SetCreatedAt(old).
		SetScrapeStatus(entlink.ScrapeStatusSuccess).SetScrapeAttemptedAt(old).SetScrapedAt(old).
		Save(context.Background())
	require.NoError(t, err)
	return row
}

func TestLivenessPreservesMetadataAndRecovers(t *testing.T) {
	db, c := newTestChecker(t)
	ctx := context.Background()
	row := storedLink(t, db, "https://example.com/old")
	owner, err := db.Account.Create().SetHandle("owner").SetName("Owner").Save(ctx)
	require.NoError(t, err)
	image, err := db.Asset.Create().SetAccountID(owner.ID).SetFilename("saved.png").SetSize(1).SetMimeType("image/png").Save(ctx)
	require.NoError(t, err)
	require.NoError(t, db.Link.UpdateOneID(row.ID).SetPrimaryImageID(image.ID).SetFaviconImageID(image.ID).Exec(ctx))
	c.prober = probeFunc(func(context.Context, link_url.URL) (link_check.Status, error) {
		return link_check.Unknown, errors.New("connection refused")
	})

	for attempt := range 3 {
		require.NoError(t, c.check(ctx, link_check.Candidate{ID: row.ID, URL: row.URL, Failures: attempt}))
		stored, err := db.Link.Get(ctx, row.ID)
		require.NoError(t, err)
		if attempt < 2 {
			require.Equal(t, entlink.LivenessStatusUnknown, stored.LivenessStatus)
		} else {
			require.Equal(t, entlink.LivenessStatusOffline, stored.LivenessStatus)
		}
		require.Equal(t, attempt+1, stored.LivenessFailures)
		require.NotNil(t, stored.LivenessCheckedAt)
		require.Equal(t, "Saved title", stored.Title)
		require.Equal(t, "Saved description", stored.Description)
		require.Equal(t, entlink.ScrapeStatusSuccess, stored.ScrapeStatus)
		require.WithinDuration(t, *row.ScrapedAt, *stored.ScrapedAt, time.Microsecond)
		require.WithinDuration(t, *row.ScrapeAttemptedAt, *stored.ScrapeAttemptedAt, time.Microsecond)
		require.Equal(t, image.ID, *stored.PrimaryAssetID)
		require.Equal(t, image.ID, *stored.FaviconAssetID)
	}

	u, err := link_url.Parse(row.URL)
	require.NoError(t, err)
	_, err = link_writer.New(db).Store(ctx, u, "", "")
	require.NoError(t, err)
	stored, err := db.Link.Get(ctx, row.ID)
	require.NoError(t, err)
	require.Equal(t, entlink.LivenessStatusOffline, stored.LivenessStatus, "saving a bookmark preserves its health")
	require.Equal(t, 3, stored.LivenessFailures)
	require.NotNil(t, stored.LivenessCheckedAt)

	c.prober = probeFunc(func(context.Context, link_url.URL) (link_check.Status, error) { return link_check.Alive, nil })
	require.NoError(t, c.check(ctx, link_check.Candidate{ID: row.ID, URL: row.URL, Failures: 3}))
	stored, err = db.Link.Get(ctx, row.ID)
	require.NoError(t, err)
	require.Equal(t, entlink.LivenessStatusAlive, stored.LivenessStatus)
	require.Zero(t, stored.LivenessFailures)
}

func TestPolicyCleanupPreservesNodesAndThreads(t *testing.T) {
	db, c := newTestChecker(t)
	ctx := context.Background()
	row := storedLink(t, db, "http://127.0.0.1/legacy")
	owner, err := db.Account.Create().SetHandle("owner").SetName("Owner").Save(ctx)
	require.NoError(t, err)
	node, err := db.Node.Create().SetName("Saved node").SetSlug("saved-node").SetAccountID(owner.ID).SetLinkID(row.ID).Save(ctx)
	require.NoError(t, err)
	post, err := db.Post.Create().SetTitle("Saved thread").SetSlug("saved-thread").SetAccountPosts(owner.ID).
		SetUpdatedAt(time.Now()).SetLastReplyAt(time.Now()).SetBody("Saved body").SetShort("Saved body").SetLinkID(row.ID).Save(ctx)
	require.NoError(t, err)
	c.prober = probeFunc(func(context.Context, link_url.URL) (link_check.Status, error) {
		t.Fatal("a policy-invalid URL must not be requested")
		return link_check.Unknown, nil
	})

	require.NoError(t, c.check(ctx, link_check.Candidate{ID: row.ID, URL: row.URL}))
	exists, err := db.Link.Query().Where(entlink.ID(row.ID)).Exist(ctx)
	require.NoError(t, err)
	require.False(t, exists)
	keptNode, err := db.Node.Get(ctx, node.ID)
	require.NoError(t, err)
	require.Equal(t, "Saved node", keptNode.Name)
	require.Equal(t, xid.NilID(), keptNode.LinkID)
	keptPost, err := db.Post.Get(ctx, post.ID)
	require.NoError(t, err)
	require.Equal(t, "Saved body", keptPost.Body)
	require.Equal(t, xid.NilID(), keptPost.LinkID)
}

func TestNetworkPolicyFailureIsRemovedAndCancellationIsNotRecorded(t *testing.T) {
	db, c := newTestChecker(t)
	ctx := context.Background()
	row := storedLink(t, db, "https://example.com/unsafe-redirect")
	c.prober = probeFunc(func(context.Context, link_url.URL) (link_check.Status, error) {
		return link_check.Unknown, fault.New("blocked redirect", ftag.With(ftag.InvalidArgument))
	})
	require.NoError(t, c.check(ctx, link_check.Candidate{ID: row.ID, URL: row.URL}))
	_, err := db.Link.Get(ctx, row.ID)
	require.True(t, ent.IsNotFound(err))

	row = storedLink(t, db, "https://example.com/cancelled")
	run, cancel := context.WithCancel(ctx)
	defer cancel()
	c.prober = probeFunc(func(ctx context.Context, _ link_url.URL) (link_check.Status, error) {
		cancel()
		return link_check.Unknown, ctx.Err()
	})
	require.ErrorIs(t, c.check(run, link_check.Candidate{ID: row.ID, URL: row.URL}), context.Canceled)
	stored, err := db.Link.Get(ctx, row.ID)
	require.NoError(t, err)
	require.Nil(t, stored.LivenessCheckedAt)
	require.Zero(t, stored.LivenessFailures)
}

func TestSweepUsesAllocatedBatchAndSkipsFreshChecks(t *testing.T) {
	db, c := newTestChecker(t)
	ctx := context.Background()
	for range maxBatchSize + 2 {
		storedLink(t, db, "https://example.com/"+xid.New().String())
	}
	_, err := db.Link.Create().SetURL("https://example.com/new").SetSlug("fresh").SetDomain("example.com").
		SetTitle("New link").SetDescription("").Save(ctx)
	require.NoError(t, err)
	recent := storedLink(t, db, "https://example.com/recent")
	require.NoError(t, db.Link.UpdateOneID(recent.ID).SetLivenessCheckedAt(time.Now()).Exec(ctx))
	calls := 0
	c.prober = probeFunc(func(context.Context, link_url.URL) (link_check.Status, error) {
		calls++
		return link_check.Alive, nil
	})
	_, err = c.run(ctx, &jobs.Job[Payload]{Payload: Payload{Limit: 1000}})
	require.NoError(t, err)
	require.Equal(t, 100, calls)
	_, err = c.run(ctx, &jobs.Job[Payload]{Payload: Payload{Limit: 1}})
	require.NoError(t, err)
	require.Equal(t, 101, calls)
	remaining, err := c.links.Due(ctx, time.Now(), 10)
	require.NoError(t, err)
	require.Len(t, remaining, 1)
}

func TestDueRetriesUncertainFailuresBeforeOfflineLinks(t *testing.T) {
	db, c := newTestChecker(t)
	ctx := context.Background()
	now := time.Now()
	uncertain := storedLink(t, db, "https://example.com/uncertain")
	offline := storedLink(t, db, "https://example.com/offline")
	unchecked := storedLink(t, db, "https://example.com/unchecked")
	oldAlive := storedLink(t, db, "https://example.com/old-alive")
	oldOffline := storedLink(t, db, "https://example.com/old-offline")
	for _, row := range []*ent.Link{uncertain, offline} {
		require.NoError(t, db.Link.UpdateOneID(row.ID).SetLivenessCheckedAt(now.Add(-link_check.RetryInterval-time.Hour)).Exec(ctx))
	}
	require.NoError(t, db.Link.UpdateOneID(uncertain.ID).SetLivenessFailures(1).Exec(ctx))
	require.NoError(t, db.Link.UpdateOneID(offline.ID).SetLivenessFailures(3).Exec(ctx))
	require.NoError(t, db.Link.UpdateOneID(oldAlive.ID).SetLivenessCheckedAt(now.Add(-link_check.CheckInterval-time.Hour)).
		SetLivenessStatus(entlink.LivenessStatusAlive).Exec(ctx))
	require.NoError(t, db.Link.UpdateOneID(oldOffline.ID).SetLivenessCheckedAt(now.Add(-link_check.CheckInterval-24*time.Hour)).
		SetLivenessStatus(entlink.LivenessStatusOffline).SetLivenessFailures(3).Exec(ctx))
	due, err := c.links.Due(ctx, now, 10)
	require.NoError(t, err)
	require.Len(t, due, 4)
	require.Equal(t, unchecked.ID, due[0].ID)
	require.Equal(t, oldOffline.ID, due[1].ID)
	require.Equal(t, oldAlive.ID, due[2].ID)
	require.Equal(t, uncertain.ID, due[3].ID)

	limited, err := c.links.Due(ctx, now, 2)
	require.NoError(t, err)
	require.Equal(t, due[:2], limited)

	empty, err := c.links.Due(ctx, now, 0)
	require.NoError(t, err)
	require.Empty(t, empty)
}

func TestScheduleEnqueueIsSharedPerTimeWindow(t *testing.T) {
	db, c := newTestChecker(t)
	ctx := context.Background()
	now := time.Date(2026, 11, 10, 20, 0, 0, 0, time.UTC)
	storedLink(t, db, "https://example.com/queued")
	require.NoError(t, c.enqueue(ctx, now))
	require.NoError(t, c.enqueue(ctx, now.Add(3*time.Hour)))
	count, err := db.BackgroundJob.Query().Count(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, count)
	row, err := db.BackgroundJob.Query().Only(ctx)
	require.NoError(t, err)
	require.Equal(t, "2026-11-10T20:00:00Z", row.Key)
	queued, err := jobs.Map[Payload](row)
	require.NoError(t, err)
	require.Equal(t, 1, queued.Payload.Limit)

	require.NoError(t, c.enqueue(ctx, now.Add(4*time.Hour)))
	count, err = db.BackgroundJob.Query().Count(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, count, "windows with no allocation do not enqueue jobs")
	require.NoError(t, c.enqueue(ctx, now.Add(24*time.Hour)))
	count, err = db.BackgroundJob.Query().Count(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, count)
}

func TestDailyBudgetScalesWithinTheCap(t *testing.T) {
	for _, scenario := range []struct {
		links int
		want  int
	}{
		{links: 0, want: 0},
		{links: 365, want: 1},
		{links: 366, want: 2},
		{links: 641, want: 2},
		{links: 6410, want: 18},
		{links: 64100, want: 176},
		{links: 200000, want: 548},
		{links: 1000000, want: 600},
	} {
		require.Equal(t, scenario.want, dailyBudget(scenario.links), "archive size %d", scenario.links)
	}
}

func TestBatchBudgetSpreadsTheDailyAllowanceAcrossUTCWindows(t *testing.T) {
	for _, scenario := range []struct {
		links   int
		daily   int
		batches []int
	}{
		{links: 0, daily: 0, batches: []int{0, 0, 0, 0, 0, 0}},
		{links: 365, daily: 1, batches: []int{0, 0, 0, 0, 0, 1}},
		{links: 641, daily: 2, batches: []int{0, 0, 1, 0, 0, 1}},
		{links: 6410, daily: 18, batches: []int{3, 3, 3, 3, 3, 3}},
		{links: 64100, daily: 176, batches: []int{29, 29, 30, 29, 29, 30}},
		{links: 200000, daily: 548, batches: []int{91, 91, 92, 91, 91, 92}},
		{links: 1000000, daily: 600, batches: []int{100, 100, 100, 100, 100, 100}},
	} {
		total := 0
		for window, want := range scenario.batches {
			now := time.Date(2026, 11, 10, window*4, 0, 0, 0, time.UTC)
			got := batchBudget(scenario.links, now)
			require.Equal(t, want, got, "archive size %d, window %d", scenario.links, window)
			require.LessOrEqual(t, got, 100)
			require.Equal(t, got, batchBudget(scenario.links, now.Add(3*time.Hour+59*time.Minute)))
			require.Equal(t, got, batchBudget(scenario.links, now.In(time.FixedZone("UTC+14", 14*60*60))))
			total += got
		}
		require.Equal(t, scenario.daily, total, "archive size %d", scenario.links)
	}
}

func TestScheduleSkipsEmptyAndNotDueArchives(t *testing.T) {
	db, c := newTestChecker(t)
	ctx := context.Background()
	now := time.Date(2026, 11, 10, 20, 0, 0, 0, time.UTC)
	require.NoError(t, c.enqueue(ctx, now))
	row := storedLink(t, db, "https://example.com/recent-check")
	require.NoError(t, db.Link.UpdateOneID(row.ID).SetLivenessCheckedAt(now).Exec(ctx))
	require.NoError(t, c.enqueue(ctx, now))
	count, err := db.BackgroundJob.Query().Count(ctx)
	require.NoError(t, err)
	require.Zero(t, count)
}

func TestDueSharesCapacityAndFillsUnusedSlots(t *testing.T) {
	db, c := newTestChecker(t)
	ctx := context.Background()
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	unchecked := make(map[xid.ID]bool)
	checked := make(map[xid.ID]bool)
	for range 6 {
		row := storedLink(t, db, "https://example.com/new-"+xid.New().String())
		unchecked[row.ID] = true
		row = storedLink(t, db, "https://example.com/old-"+xid.New().String())
		checked[row.ID] = true
		require.NoError(t, db.Link.UpdateOneID(row.ID).
			SetLivenessCheckedAt(now.Add(-link_check.CheckInterval-time.Hour)).Exec(ctx))
	}

	due, err := c.links.Due(ctx, now, 4)
	require.NoError(t, err)
	require.Len(t, due, 4)
	uncheckedCount := 0
	for _, candidate := range due {
		if unchecked[candidate.ID] {
			uncheckedCount++
		}
	}
	require.Equal(t, 2, uncheckedCount)

	first, err := c.links.Due(ctx, now, 1)
	require.NoError(t, err)
	second, err := c.links.Due(ctx, now.Add(24*time.Hour), 1)
	require.NoError(t, err)
	require.Len(t, first, 1)
	require.Len(t, second, 1)
	require.NotEqual(t, unchecked[first[0].ID], unchecked[second[0].ID])

	for id := range checked {
		require.NoError(t, db.Link.UpdateOneID(id).SetLivenessCheckedAt(now).Exec(ctx))
	}
	due, err = c.links.Due(ctx, now, 4)
	require.NoError(t, err)
	require.Len(t, due, 4)
	for _, candidate := range due {
		require.True(t, unchecked[candidate.ID])
	}

	for id := range unchecked {
		require.NoError(t, db.Link.UpdateOneID(id).
			SetLivenessCheckedAt(now.Add(-link_check.CheckInterval-time.Hour)).Exec(ctx))
	}
	due, err = c.links.Due(ctx, now, 4)
	require.NoError(t, err)
	require.Len(t, due, 4)
	for _, candidate := range due {
		require.True(t, unchecked[candidate.ID])
	}
}

func TestDueUsesArchiveTimescales(t *testing.T) {
	db, c := newTestChecker(t)
	ctx := context.Background()
	now := time.Now().UTC()
	want := make(map[xid.ID]bool)
	for _, scenario := range []struct {
		name     string
		age      time.Duration
		checked  time.Duration
		failures int
		due      bool
	}{
		{name: "young", age: 89 * 24 * time.Hour},
		{name: "aged", age: 91 * 24 * time.Hour, due: true},
		{name: "healthy-recent", age: 730 * 24 * time.Hour, checked: 364 * 24 * time.Hour},
		{name: "healthy-old", age: 730 * 24 * time.Hour, checked: 366 * 24 * time.Hour, due: true},
		{name: "failure-recent", age: 730 * 24 * time.Hour, checked: 29 * 24 * time.Hour, failures: 1},
		{name: "failure-retry", age: 730 * 24 * time.Hour, checked: 31 * 24 * time.Hour, failures: 2, due: true},
		{name: "offline-recent", age: 730 * 24 * time.Hour, checked: 31 * 24 * time.Hour, failures: 3},
		{name: "offline-old", age: 730 * 24 * time.Hour, checked: 366 * 24 * time.Hour, failures: 3, due: true},
	} {
		create := db.Link.Create().SetURL("https://example.com/" + scenario.name).
			SetSlug(scenario.name).SetDomain("example.com").SetTitle(scenario.name).SetDescription("").
			SetCreatedAt(now.Add(-scenario.age)).SetLivenessFailures(scenario.failures)
		if scenario.checked > 0 {
			create.SetLivenessCheckedAt(now.Add(-scenario.checked))
		}
		row, err := create.Save(ctx)
		require.NoError(t, err)
		if scenario.due {
			want[row.ID] = true
		}
	}

	due, err := c.links.Due(ctx, now, 20)
	require.NoError(t, err)
	require.Len(t, due, len(want))
	for _, candidate := range due {
		require.True(t, want[candidate.ID])
	}
}
