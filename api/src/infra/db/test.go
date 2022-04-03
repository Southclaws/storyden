package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/Southclaws/storyden/api/src/infra/db/model"
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
		// tables := []string{
		// 	notification.Table,
		// 	subscription.Table,
		// 	react.Table,
		// 	rule.Table,
		// 	server.Table,
		// 	user.Table,
		// 	category.Table,
		// 	tag.Table,
		// 	post.Table,
		// }

		// q := fmt.Sprintf("truncate table %s CASCADE;", strings.Join(tables, ", "))

		// if _, err := d.Exec(q); err != nil {
		// 	t.Fatal(err)
		// }
		c.Close()

		fmt.Println("--- Cleaned database after test", d)
	})

	return c
}
