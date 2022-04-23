package thread

import (
	"fmt"

	"github.com/Southclaws/storyden/api/src/infra/db/model"
)

func NewLocalWithSeed() Repository {
	m := NewLocal()
	Seed(m)
	return m
}

func NewWithSeed(db *model.Client) Repository {
	m := New(db)
	Seed(m)
	return m
}

func Seed(r Repository) {
	// ctx := context.Background()

	fmt.Println("created seed threads")
}
