package main

import (
	"github.com/Southclaws/storyden/internal/script"
	"github.com/Southclaws/storyden/pkg/resources/seed"
)

func main() {
	script.Run(seed.Create())
}
