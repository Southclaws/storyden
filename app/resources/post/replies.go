package post

import (
	"github.com/rs/xid"
	"github.com/samber/lo"
)

type ReplyStatus struct {
	Count   int
	Replied int
}

type PostRepliesResult struct {
	PostID  xid.ID `db:"post_id"`
	Count   int    `db:"replies"`
	Replied int    `db:"replied"`
}

func (p PostRepliesResult) Status() ReplyStatus {
	return ReplyStatus{
		Count:   p.Count,
		Replied: p.Replied,
	}
}

type PostRepliesResults []PostRepliesResult

func (p PostRepliesResults) Map() PostRepliesMap {
	return lo.KeyBy(p, func(x PostRepliesResult) xid.ID { return x.PostID })
}

// PostRepliesMap maps the likes aggregation back to individual posts.
type PostRepliesMap map[xid.ID]PostRepliesResult

func (p PostRepliesMap) Status(id xid.ID) ReplyStatus {
	if p == nil {
		return ReplyStatus{}
	}

	s, ok := p[id]
	if !ok {
		return ReplyStatus{}
	}

	return s.Status()
}
