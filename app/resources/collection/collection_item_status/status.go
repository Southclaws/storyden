package collection_item_status

import (
	"github.com/rs/xid"
	"github.com/samber/lo"
)

type Status struct {
	Count  int
	Status bool
}

type CollectionStatusResult struct {
	ItemID    xid.ID `db:"item_id"`
	Count     int    `db:"collections"`
	Collected int    `db:"has_in_collection"`
}

func (p CollectionStatusResult) Status() Status {
	return Status{
		Count:  p.Count,
		Status: p.Collected > 0,
	}
}

type CollectionStatusResults []CollectionStatusResult

func (p CollectionStatusResults) Map() CollectionStatusMap {
	return lo.KeyBy(p, func(x CollectionStatusResult) xid.ID { return x.ItemID })
}

type CollectionStatusMap map[xid.ID]CollectionStatusResult

func (p CollectionStatusMap) Status(id xid.ID) Status {
	if p == nil {
		return Status{}
	}

	s, ok := p[id]
	if !ok {
		return Status{}
	}

	return s.Status()
}
