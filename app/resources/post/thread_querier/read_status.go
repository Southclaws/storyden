package thread_querier

import (
	"context"
	"fmt"
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/post"
)

const newRepliesCountManyQuery_sqlite = `select
  p.id        as post_id,
  max(pr.last_seen_at) as last_read_at,
  count(r.id) as new_replies
from
  posts p
  inner join post_reads pr
    on pr.root_post_id = p.id and pr.account_id = $1
  left join posts r
    on r.root_post_id = p.id
    and r.deleted_at is null
    and unixepoch(r.created_at) > unixepoch(pr.last_seen_at)
where p.id in (%s)
group by p.id
`

const newRepliesCountManyQuery_postgres = `select
  p.id        as post_id,
  max(pr.last_seen_at) as last_read_at,
  count(r.id) as new_replies
from
  posts p
  inner join post_reads pr
    on pr.root_post_id = p.id and pr.account_id = $1
  left join posts r
    on r.root_post_id = p.id
    and r.deleted_at is null
    and r.created_at > pr.last_seen_at
where p.id in (%s)
group by p.id
`

func (d *Querier) newRepliesCountManyQuery() string {
	switch d.raw.DriverName() {
	case "sqlite", "sqlite3", "libsql":
		return newRepliesCountManyQuery_sqlite
	case "pgx", "postgres":
		return newRepliesCountManyQuery_postgres
	default:
		return newRepliesCountManyQuery_postgres
	}
}

func (d *Querier) getReadStatus(ctx context.Context, ids []xid.ID, accountID string) (post.ReadStateMap, error) {
	if len(ids) == 0 {
		return post.ReadStateMap{}, nil
	}

	quotedIDs := dt.Map(ids, func(id xid.ID) string { return fmt.Sprintf("'%s'", id.String()) })
	idList := strings.Join(quotedIDs, ",")

	var readStates post.ReadStateResults
	readQuery := fmt.Sprintf(d.newRepliesCountManyQuery(), idList)
	err := d.raw.SelectContext(ctx, &readStates, readQuery, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return readStates.Map(), nil
}
