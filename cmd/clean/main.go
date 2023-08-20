package main

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/internal/script"
)

func main() {
	script.Run(fx.Invoke(seed.NewEmpty))
}
