package bindings

import (
	"context"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/pkg/transports/http/openapi"
)

type Version struct{}

func NewVersion() Version { return Version{} }

func (v *Version) GetVersion(ctx context.Context, request openapi.GetVersionRequestObject) any {
	return openapi.GetVersion200TextResponse(config.Version)
}
