package node_traversal

import (
	"testing"

	"github.com/Southclaws/opt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/library"
)

func TestBuildSubtreeQueryUsesDriverNeutralPlaceholders(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	depth := uint(10)
	query, args := buildSubtreeQuery(opt.NewEmpty[library.NodeID](), filters{depth: &depth}, func(q string) string {
		return q
	})

	r.Contains(query, "depth <= ?")
	r.NotContains(query, "$1")
	r.Equal([]interface{}{depth}, args)
}

func TestBuildSubtreeQueryRebindsPostgresPlaceholders(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	nodeID := library.NodeID(xid.New())
	depth := uint(10)
	handle := "southclaws"
	query, args := buildSubtreeQuery(opt.New(nodeID), filters{
		rootAccountHandleFilter: &handle,
		depth:                   &depth,
	}, func(q string) string {
		return sqlx.Rebind(sqlx.DOLLAR, q)
	})

	r.Contains(query, "id = cast($1 as text)")
	r.Contains(query, "a.handle = $2")
	r.Contains(query, "depth <= $3")
	r.Equal([]interface{}{nodeID.String(), handle, depth}, args)
}

func TestSQLXRebindsStorydenDriverNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		driver string
		want   string
	}{
		{name: "postgres pgx", driver: "pgx", want: "depth <= $1"},
		{name: "cockroach", driver: "cockroach", want: "depth <= $1"},
		{name: "sqlite", driver: "sqlite", want: "depth <= ?"},
		{name: "libsql", driver: "libsql", want: "depth <= ?"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, sqlx.Rebind(sqlx.BindType(tc.driver), "depth <= ?"))
		})
	}
}
