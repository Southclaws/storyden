package thread_querier

import (
	"context"
	"fmt"
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/collection/collection_item_status"
)

const collectionsCountManyQuery = `select
  p.id        item_id,
  count(cp.post_id) collections,
  count(a.id) has_in_collection
from
  posts p
  left join collection_posts cp on cp.post_id = p.id
  left join collections c on c.id = cp.collection_id
  left join accounts a on c.account_collections = a.id and a.id = $1
where p.id in (%s)
group by p.id
`

func (d *Querier) getCollectionsStatus(ctx context.Context, ids []xid.ID, accountID string) (collection_item_status.CollectionStatusMap, error) {
	if len(ids) == 0 {
		return collection_item_status.CollectionStatusMap{}, nil
	}

	quotedIDs := dt.Map(ids, func(id xid.ID) string { return fmt.Sprintf("'%s'", id.String()) })
	idList := strings.Join(quotedIDs, ",")

	var collections collection_item_status.CollectionStatusResults
	collectionsQuery := fmt.Sprintf(collectionsCountManyQuery, idList)
	err := d.raw.SelectContext(ctx, &collections, collectionsQuery, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return collections.Map(), nil
}
