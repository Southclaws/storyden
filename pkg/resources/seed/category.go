package seed

import (
	"context"
	"fmt"

	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/pkg/resources/category"
)

var (
	Category_01_General = category.Category{
		ID:          category.CategoryID(id("00000000000000000010")),
		Name:        "General",
		Description: "General stuff",
		Colour:      "#ffffff",
		Sort:        0,
	}

	Category_02_Photos = category.Category{
		ID:          category.CategoryID(id("00000000000000000020")),
		Name:        "Media",
		Description: "Movies and tv shows",
		Colour:      "#ffffff",
		Sort:        1,
	}

	Category_03_Movies = category.Category{
		ID:          category.CategoryID(id("00000000000000000030")),
		Name:        "Movies",
		Description: "Movies discussion",
		Colour:      "#ffffff",
		Sort:        2,
	}

	Category_04_Music = category.Category{
		ID:          category.CategoryID(id("00000000000000000040")),
		Name:        "Music",
		Description: "Music discussion",
		Colour:      "#ffffff",
		Sort:        3,
	}

	Category_05_Admin = category.Category{
		ID:          category.CategoryID(id("00000000000000000050")),
		Name:        "Admin",
		Description: "Admin area",
		Colour:      "#ffffff",
		Sort:        4,
		Admin:       true,
	}
)

func categories(r category.Repository) {
	ctx := context.Background()

	for _, c := range []category.Category{
		Category_01_General,
		Category_02_Photos,
		Category_03_Movies,
		Category_04_Music,
		Category_05_Admin,
	} {
		utils.Must(r.CreateCategory(ctx, c.Name, c.Description, c.Colour, c.Sort, c.Admin, category.WithID(c.ID)))
	}

	fmt.Println("created seed categories")
}
