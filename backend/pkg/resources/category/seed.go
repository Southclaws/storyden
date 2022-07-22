package category

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
)

var (
	SeedCategory_01_General = Category{
		ID:          CategoryID(uuid.MustParse("00000000-0000-0000-0000-000000000000")),
		Name:        "General",
		Description: "General stuff",
		Colour:      "#ffffff",
		Sort:        0,
	}

	SeedCategory_02_Photos = Category{
		ID:          CategoryID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		Name:        "Media",
		Description: "Movies and tv shows",
		Colour:      "#ffffff",
		Sort:        1,
	}

	SeedCategory_03_Movies = Category{
		ID:          CategoryID(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
		Name:        "Movies",
		Description: "Movies discussion",
		Colour:      "#ffffff",
		Sort:        2,
	}

	SeedCategory_04_Music = Category{
		ID:          CategoryID(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
		Name:        "Music",
		Description: "Music discussion",
		Colour:      "#ffffff",
		Sort:        3,
	}

	SeedCategory_05_Admin = Category{
		ID:          CategoryID(uuid.MustParse("00000000-0000-0000-0000-000000000004")),
		Name:        "Admin",
		Description: "Admin area",
		Colour:      "#ffffff",
		Sort:        4,
		Admin:       true,
	}
)

// func NewLocalWithSeed() Repository {
// 	m := NewLocal()
// 	Seed(m)
// 	return m
// }

func NewWithSeed(db *model.Client) Repository {
	m := New(db)
	Seed(m)
	return m
}

func Seed(r Repository) {
	ctx := context.Background()

	create := func(c *Category) CategoryID {
		c, err := r.CreateCategory(ctx, c.Name, c.Description, c.Colour, c.Sort, c.Admin)
		if err != nil {
			panic(err)
		}
		return c.ID
	}

	SeedCategory_01_General.ID = create(&SeedCategory_01_General)
	SeedCategory_02_Photos.ID = create(&SeedCategory_02_Photos)
	SeedCategory_03_Movies.ID = create(&SeedCategory_03_Movies)
	SeedCategory_04_Music.ID = create(&SeedCategory_04_Music)
	SeedCategory_05_Admin.ID = create(&SeedCategory_05_Admin)

	fmt.Println("created seed categories")
}
