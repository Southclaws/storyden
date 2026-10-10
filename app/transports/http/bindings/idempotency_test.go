package bindings

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	_ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/idempotency"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/enttest"
)

func TestIdempotentCreateRechecksReplayVisibility(t *testing.T) {
	raw, err := sql.Open("sqlite", "file:"+filepath.Join(t.TempDir(), "receipts.db")+"?_pragma=foreign_keys(1)")
	require.NoError(t, err)
	db := enttest.NewClient(t, enttest.WithOptions(ent.Driver(entsql.OpenDB(dialect.SQLite, raw))))
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	receipts := idempotency.New(db)
	key := openapi.IdempotencyKey("create-one")
	creates := 0
	visible := true
	create := func() (struct{ ID string }, error) {
		creates++
		return struct{ ID string }{ID: "resource-one"}, nil
	}
	recheck := func(_ struct{ ID string }) error {
		if !visible {
			return errors.New("resource is no longer visible")
		}
		return nil
	}
	ctx := context.Background()
	first, err := idempotentCreate(ctx, receipts, "account", "ThreadCreate", &key, map[string]string{"title": "A"}, create, recheck)
	require.NoError(t, err)
	require.Equal(t, "resource-one", first.ID)
	_, err = idempotentCreate(ctx, receipts, "account", "ThreadCreate", &key, map[string]string{"title": "B"}, create, recheck)
	require.Error(t, err)
	require.Equal(t, ftag.AlreadyExists, ftag.Get(err))
	require.Equal(t, "key_mismatch", fctx.Unwrap(err)["idempotency_reason"])
	require.Equal(t, 1, creates)
	visible = false
	_, err = idempotentCreate(ctx, receipts, "account", "ThreadCreate", &key, map[string]string{"title": "A"}, create, recheck)
	require.ErrorContains(t, err, "no longer visible")
	require.Equal(t, 1, creates)
}
