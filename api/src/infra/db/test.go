package db

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Southclaws/storyden/api/src/infra/db/model"
	"github.com/Southclaws/storyden/api/src/infra/db/model/category"
	"github.com/Southclaws/storyden/api/src/infra/db/model/notification"
	"github.com/Southclaws/storyden/api/src/infra/db/model/post"
	"github.com/Southclaws/storyden/api/src/infra/db/model/react"
	"github.com/Southclaws/storyden/api/src/infra/db/model/rule"
	"github.com/Southclaws/storyden/api/src/infra/db/model/server"
	"github.com/Southclaws/storyden/api/src/infra/db/model/subscription"
	"github.com/Southclaws/storyden/api/src/infra/db/model/tag"
	"github.com/Southclaws/storyden/api/src/infra/db/model/user"
)

func TestDB(t *testing.T) *model.Client {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgresql://default:default@localhost:5432/postgres"
	}

	c, d, err := connect(url)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		tables := []string{
			notification.Table,
			subscription.Table,
			react.Table,
			rule.Table,
			server.Table,
			user.Table,
			category.Table,
			tag.Table,
			post.Table,
		}

		q := fmt.Sprintf("truncate table %s CASCADE;", strings.Join(tables, ", "))

		if _, err := d.Exec(q); err != nil {
			t.Fatal(err)
		}
		c.Close()

		fmt.Println("--- Cleaned database after test")
	})

	return c
}
