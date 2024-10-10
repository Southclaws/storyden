package library

import (
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/google/uuid"
)

func UniqueNode(name string) openapi.NodeInitialProps {
	slug := name + uuid.NewString()
	return openapi.NodeInitialProps{
		Name: name,
		Slug: &slug,
	}
}
