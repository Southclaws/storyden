package storyden

//go:build+ tools

import (
	_ "entgo.io/ent"
	_ "github.com/99designs/gqlgen/codegen"
	_ "github.com/deepmap/oapi-codegen/pkg/runtime"
)
