package idempotency

import (
	"context"
	"database/sql"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/enttest"
)

func testRepositories(t *testing.T) (*Repository, *Repository, *sql.DB) {
	t.Helper()
	dsn := "file:" + filepath.Join(t.TempDir(), "receipts.db") + "?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"
	var firstRaw *sql.DB
	open := func() *ent.Client {
		raw, err := sql.Open("sqlite", dsn)
		require.NoError(t, err)
		if firstRaw == nil {
			firstRaw = raw
		}
		raw.SetMaxOpenConns(1)
		client := enttest.NewClient(t, enttest.WithOptions(ent.Driver(entsql.OpenDB(dialect.SQLite, raw))))
		t.Cleanup(func() { require.NoError(t, client.Close()) })
		return client
	}
	return New(open()), New(open()), firstRaw
}

func TestClaimAcrossDatabaseClients(t *testing.T) {
	a, b, _ := testRepositories(t)
	ctx := context.Background()
	results := make([]Result, 2)
	errors := make([]error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		results[0], errors[0] = a.Claim(ctx, "account-a", "ThreadCreate", "key", "body-a")
	}()
	go func() {
		defer wg.Done()
		results[1], errors[1] = b.Claim(ctx, "account-a", "ThreadCreate", "key", "body-a")
	}()
	wg.Wait()
	require.NoError(t, errors[0])
	require.NoError(t, errors[1])
	require.NotEqual(t, results[0].State, results[1].State)
	require.Contains(t, []State{results[0].State, results[1].State}, Claimed)
	require.Contains(t, []State{results[0].State, results[1].State}, InProgress)

	claim := results[0]
	if claim.State != Claimed {
		claim = results[1]
	}
	require.NoError(t, a.Complete(ctx, claim.ID, claim.Token, `{"id":"created"}`))
	replay, err := b.Claim(ctx, "account-a", "ThreadCreate", "key", "body-a")
	require.NoError(t, err)
	require.Equal(t, Completed, replay.State)
	require.Equal(t, `{"id":"created"}`, replay.Response)

	mismatch, err := b.Claim(ctx, "account-a", "ThreadCreate", "key", "body-b")
	require.NoError(t, err)
	require.Equal(t, Mismatch, mismatch.State)

	otherPrincipal, err := b.Claim(ctx, "account-b", "ThreadCreate", "key", "body-a")
	require.NoError(t, err)
	require.Equal(t, Claimed, otherPrincipal.State)
	otherOperation, err := b.Claim(ctx, "account-a", "ReplyCreate", "key", "body-a")
	require.NoError(t, err)
	require.Equal(t, Claimed, otherOperation.State)
	otherKey, err := b.Claim(ctx, "account-a", "ThreadCreate", "other", "body-a")
	require.NoError(t, err)
	require.Equal(t, Claimed, otherKey.State)
}

func TestUncertainClaimDoesNotReexecute(t *testing.T) {
	a, b, raw := testRepositories(t)
	ctx := context.Background()
	claim, err := a.Claim(ctx, "account", "NodeCreate", "key", "body")
	require.NoError(t, err)
	require.Equal(t, Claimed, claim.State)

	_, err = raw.ExecContext(ctx, "UPDATE idempotency_receipts SET expires_at = ? WHERE id = ?", time.Now().Add(-time.Minute), claim.ID)
	require.NoError(t, err)
	retry, err := b.Claim(ctx, "account", "NodeCreate", "key", "body")
	require.NoError(t, err)
	require.Equal(t, Uncertain, retry.State)

	failed, err := a.Claim(ctx, "account", "ReplyCreate", "failed", "body")
	require.NoError(t, err)
	require.NoError(t, a.Fail(ctx, failed.ID, failed.Token))
	retry, err = b.Claim(ctx, "account", "ReplyCreate", "failed", "body")
	require.NoError(t, err)
	require.Equal(t, Uncertain, retry.State)
}

func TestExpiredCompletedClaimCanBeReusedButOldWorkerCannotCompleteNewClaim(t *testing.T) {
	a, b, raw := testRepositories(t)
	ctx := context.Background()
	first, err := a.Claim(ctx, "account", "ThreadCreate", "key", "original")
	require.NoError(t, err)
	require.NoError(t, a.Complete(ctx, first.ID, first.Token, `{"id":"first"}`))
	_, err = raw.ExecContext(ctx, "UPDATE idempotency_receipts SET expires_at = ? WHERE id = ?", time.Now().Add(-time.Minute), first.ID)
	require.NoError(t, err)

	second, err := b.Claim(ctx, "account", "ThreadCreate", "key", "changed")
	require.NoError(t, err)
	require.Equal(t, Claimed, second.State)
	require.NotEqual(t, first.Token, second.Token)
	require.Error(t, a.Complete(ctx, first.ID, first.Token, `{"id":"stale"}`))
	require.NoError(t, b.Complete(ctx, second.ID, second.Token, `{"id":"second"}`))
	got, err := a.Claim(ctx, "account", "ThreadCreate", "key", "changed")
	require.NoError(t, err)
	require.Equal(t, Completed, got.State)
	require.Equal(t, `{"id":"second"}`, got.Response)
}
