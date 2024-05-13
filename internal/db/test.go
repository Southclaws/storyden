package db

import (
	"database/sql"
	"fmt"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/asset"
	"github.com/Southclaws/storyden/internal/ent/authentication"
	"github.com/Southclaws/storyden/internal/ent/category"
	"github.com/Southclaws/storyden/internal/ent/cluster"
	"github.com/Southclaws/storyden/internal/ent/collection"
	"github.com/Southclaws/storyden/internal/ent/link"
	"github.com/Southclaws/storyden/internal/ent/notification"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/react"
	"github.com/Southclaws/storyden/internal/ent/role"
	"github.com/Southclaws/storyden/internal/ent/setting"
	"github.com/Southclaws/storyden/internal/ent/tag"
)

func Truncate(db *sql.DB) error {
	tables := []string{
		tag.Table,
		setting.Table,
		role.Table,
		react.Table,
		post.Table,
		notification.Table,
		link.Table,
		collection.Table,
		cluster.Table,
		category.Table,
		authentication.Table,
		asset.Table,
		account.Table,
	}

	for _, t := range tables {
		if _, err := db.Exec(fmt.Sprintf("delete from %s", t)); err != nil {
			return fault.Wrap(err, fmsg.With(fmt.Sprintf("failed to clean table %s", t)))
		}
	}

	fmt.Println("--- Cleaned database")

	return nil
}
