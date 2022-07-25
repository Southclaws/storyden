package db

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/category"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/notification"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/post"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/react"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/subscription"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/tag"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/user"
)

func TestDB(t *testing.T) *model.Client {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgresql://default:default@localhost:5432/postgres"
	}

	c, d, err := connect(context.Background(), url, false)
	if err != nil {
		t.Fatal(err)
	}

	truncate := func() {
		tables := []string{
			notification.Table,
			subscription.Table,
			react.Table,
			user.Table,
			category.Table,
			tag.Table,
			post.Table,
		}

		q := fmt.Sprintf("truncate table %s CASCADE;", strings.Join(tables, ", "))

		if _, err := d.Exec(q); err != nil {
			t.Fatal(err)
		}
	}

	// truncate the database before tests and after.
	truncate()
	t.Cleanup(func() {
		truncate()
		c.Close()

		fmt.Println("--- Cleaned database after test")
	})

	return c
}
