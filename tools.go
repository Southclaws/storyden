package storyden

//go:build+ tools

import (
	_ "entgo.io/ent/cmd/ent"
	_ "github.com/99designs/gqlgen"
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"
)
