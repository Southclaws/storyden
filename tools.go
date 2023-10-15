//go:build tools
// +build tools

package storyden

import (
	_ "entgo.io/ent"
	_ "github.com/99designs/gqlgen/codegen"
	_ "github.com/Southclaws/enumerator"
	_ "github.com/a8m/enter"
	_ "github.com/deepmap/oapi-codegen/pkg/runtime"
)
