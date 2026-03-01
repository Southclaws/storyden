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
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
)

const repliesCountManyQuery = `select
  p.id                                                  post_id,
  count(r.id)                                           replies,
  sum(case when r.account_posts = $1 then 1 else 0 end) replied
from
  posts p
  left join posts r on r.root_post_id = p.id and r.deleted_at is null
where p.id in (%s)
group by p.id
`

func (d *Querier) getRepliesStatus(ctx context.Context, ids []xid.ID, accountID string) (post.PostRepliesMap, error) {
	ctx, span := d.ins.InstrumentNamed(ctx, "replies_status_query",
		kv.Int("id_count", len(ids)),
	)
	defer span.End()

	if len(ids) == 0 {
		return post.PostRepliesMap{}, nil
	}

	quotedIDs := dt.Map(ids, func(id xid.ID) string { return fmt.Sprintf("'%s'", id.String()) })
	idList := strings.Join(quotedIDs, ",")

	var replies post.PostRepliesResults
	repliesQuery := fmt.Sprintf(repliesCountManyQuery, idList)
	err := d.raw.SelectContext(ctx, &replies, repliesQuery, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	span.Annotate(
		kv.Int("result_rows", len(replies)),
	)

	return replies.Map(), nil
}
