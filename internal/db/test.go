package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/category"
	"github.com/Southclaws/storyden/internal/ent/notification"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/react"
	"github.com/Southclaws/storyden/internal/ent/subscription"
	"github.com/Southclaws/storyden/internal/ent/tag"
)

func Truncate(db *sql.DB) error {
	tables := []string{
		notification.Table,
		subscription.Table,
		react.Table,
		account.Table,
		category.Table,
		tag.Table,
		post.Table,
	}

	if _, err := db.Exec(fmt.Sprintf("truncate table %s CASCADE;", strings.Join(tables, ", "))); err != nil {
		return fault.Wrap(err, fmsg.With("failed to clean database"))
	}

	fmt.Println("--- Cleaned database")

	return nil
}
