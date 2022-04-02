package tag

import (
	"context"

	"github.com/Southclaws/storyden/ent"
)

type DB struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &DB{db}
}

func (d *DB) GetTags(ctx context.Context, query string) ([]Tag, error) {
	var tags []Tag
	err := d.db.Prisma.Raw.QueryRaw(`
		select
			t.id, t.name, count(*) as posts
		from
			"_TagToPost" ttp
		inner join
			"Tag" t on ttp."B" = t.id
		inner join
			"Post" p on ttp."A" = p.id
		where t.name like $1
		group by t.id
		order by posts desc
	`, query+"%").Exec(ctx, &tags)
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return nil, nil
	}

	return tags, nil
}
