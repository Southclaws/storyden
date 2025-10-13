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

const likesCountManyQuery = `select
  p.id        post_id, -- the post (thread or reply) ID
  count(*)    likes,   -- number of likes
  count(a.id) liked    -- has the account making the query liked this post?
from
  like_posts lp
  inner join posts p on p.id = lp.post_id
  left join accounts a on lp.account_id = a.id and a.id = $1
where p.id in (%s)
group by p.id
`

func (d *Querier) getLikesStatus(ctx context.Context, ids []xid.ID, accountID string) (post.PostLikesMap, error) {
	if len(ids) == 0 {
		return post.PostLikesMap{}, nil
	}

	quotedIDs := dt.Map(ids, func(id xid.ID) string { return fmt.Sprintf("'%s'", id.String()) })
	idList := strings.Join(quotedIDs, ",")

	var likes post.PostLikesResults

	// NOTE: Safe SQL parameterization for ID list. IDs are direct from a query.
	likesQuery := fmt.Sprintf(likesCountManyQuery, idList)

	err := d.raw.SelectContext(ctx, &likes, likesQuery, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return likes.Map(), nil
}
