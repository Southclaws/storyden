package main

import (
	"github.com/Southclaws/storyden/backend/internal/script"
	"github.com/Southclaws/storyden/backend/pkg/resources/seed"
)

func main() {
	script.Run(seed.Create())
}
