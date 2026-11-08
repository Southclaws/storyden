package autoload

import "github.com/Southclaws/storyden/internal/dotenv"

func init() {
	_ = dotenv.Load(".env")
}
