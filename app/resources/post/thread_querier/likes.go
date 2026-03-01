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

const likesCountManyQuery = `select
  p.id        post_id,                                      -- the post (thread or reply) ID
  count(*)    likes,                                        -- number of likes
  sum(case when lp.account_id = ? then 1 else 0 end) liked -- has the account making the query liked this post?
from
  like_posts lp
  inner join posts p on p.id = lp.post_id
where p.id in (%s)
group by p.id
`

func (d *Querier) getLikesStatus(ctx context.Context, ids []xid.ID, accountID string) (post.PostLikesMap, error) {
	ctx, span := d.ins.InstrumentNamed(ctx, "likes_status",
		kv.Int("id_count", len(ids)),
	)
	defer span.End()

	if len(ids) == 0 {
		return post.PostLikesMap{}, nil
	}

	quotedIDs := dt.Map(ids, func(id xid.ID) string { return fmt.Sprintf("'%s'", id.String()) })
	idList := strings.Join(quotedIDs, ",")

	var likes post.PostLikesResults

	// NOTE: Safe SQL parameterization for ID list. IDs are direct from a query.
	likesQuery := d.raw.Rebind(fmt.Sprintf(likesCountManyQuery, idList))

	err := d.raw.SelectContext(ctx, &likes, likesQuery, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	span.Annotate(
		kv.Int("result_rows", len(likes)),
	)

	return likes.Map(), nil
}
