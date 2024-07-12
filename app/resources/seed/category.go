package seed

import (
	"context"
	"fmt"

	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/internal/utils"
)

var (
	Category_01_General = category.Category{
		ID:          category.CategoryID(id("00000000000000000010")),
		Name:        "General",
		Description: "General stuff",
		Colour:      "#D3DFB8",
		Sort:        0,
	}

	Category_02_Photos = category.Category{
		ID:          category.CategoryID(id("00000000000000000020")),
		Name:        "Photos",
		Description: "Share your photos with the community",
		Colour:      "#F2E86D",
		Sort:        1,
	}

	Category_03_Movies = category.Category{
		ID:          category.CategoryID(id("00000000000000000030")),
		Name:        "Movies",
		Description: "Movies discussion",
		Colour:      "#C6A15B",
		Sort:        2,
	}

	Category_04_Music = category.Category{
		ID:          category.CategoryID(id("00000000000000000040")),
		Name:        "Music",
		Description: "Music, playlists and events",
		Colour:      "#A38560",
		Sort:        3,
	}

	Category_05_Admin = category.Category{
		ID:          category.CategoryID(id("00000000000000000050")),
		Name:        "Admin",
		Description: "Admin area",
		Colour:      "#574D68",
		Sort:        4,
		Admin:       true,
	}

	Categories = []category.Category{
		Category_01_General,
		Category_02_Photos,
		Category_03_Movies,
		Category_04_Music,
		Category_05_Admin,
	}
)

func categories(r category.Repository) {
	ctx := context.Background()

	for _, c := range Categories {
		utils.Must(r.CreateCategory(ctx, c.Name, c.Description, c.Colour, c.Sort, c.Admin, category.WithID(c.ID)))
	}

	fmt.Println("created seed categories")
}
