package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Southclaws/storyden/internal/infrastructure/db/model/account"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/category"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/notification"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/post"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/react"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/subscription"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/tag"
)

func TestDB(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgresql://default:default@localhost:5432/postgres"
	}

	c, d, err := connect(context.Background(), url, false)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		Truncate(d)
		c.Close()
	})
}

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
		return err
	}

	fmt.Println("--- Cleaned database")

	return nil
}
