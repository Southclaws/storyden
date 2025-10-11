package post

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"
)

type ReadStatus struct {
	Count      int
	LastReadAt time.Time
}

type ReadStateResult struct {
	PostID     xid.ID `db:"post_id"`
	NewReplies int    `db:"new_replies"`
	LastReadAt string `db:"last_read_at"`
}

func (p ReadStateResult) Status() ReadStatus {
	t, err := time.Parse(time.RFC3339Nano, p.LastReadAt)
	if err != nil {
		t, _ = time.Parse("2006-01-02 15:04:05.999999999-07:00", p.LastReadAt)
	}
	return ReadStatus{
		Count:      p.NewReplies,
		LastReadAt: t,
	}
}

type ReadStateResults []ReadStateResult

func (p ReadStateResults) Map() ReadStateMap {
	return lo.KeyBy(p, func(x ReadStateResult) xid.ID { return x.PostID })
}

// ReadStateMap maps the likes aggregation back to individual posts.
type ReadStateMap map[xid.ID]ReadStateResult

func (p ReadStateMap) Status(id xid.ID) opt.Optional[ReadStatus] {
	if p == nil {
		return opt.NewEmpty[ReadStatus]()
	}

	s, ok := p[id]
	if !ok {
		return opt.NewEmpty[ReadStatus]()
	}

	return opt.New(s.Status())
}
