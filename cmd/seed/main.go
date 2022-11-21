package main

import (
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/internal/script"
)

func main() {
	script.Run(seed.Create())
}
