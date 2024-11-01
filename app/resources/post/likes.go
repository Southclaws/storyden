package post

import (
	"github.com/Southclaws/storyden/app/resources/like"
	"github.com/rs/xid"
	"github.com/samber/lo"
)

// PostLikesResult is used for querying all members who liked posts in a thread.
type PostLikesResult struct {
	PostID xid.ID `db:"post_id"`
	Count  int    `db:"likes"`
	Liked  int    `db:"liked"`
}

func (p PostLikesResult) Status() like.Status {
	return like.Status{
		Count:  p.Count,
		Status: p.Liked > 0,
	}
}

type PostLikesResults []PostLikesResult

func (p PostLikesResults) Map() PostLikesMap {
	return lo.KeyBy(p, func(x PostLikesResult) xid.ID { return x.PostID })
}

// PostLikesMap maps the likes aggregation back to individual posts.
type PostLikesMap map[xid.ID]PostLikesResult

func (p PostLikesMap) Status(id xid.ID) like.Status {
	if p == nil {
		return like.Status{}
	}

	s, ok := p[id]
	if !ok {
		return like.Status{}
	}

	return s.Status()
}
