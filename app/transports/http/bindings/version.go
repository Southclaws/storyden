package bindings

import (
	"context"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
)

type Version struct{}

func NewVersion() Version { return Version{} }

func (v *Version) GetVersion(ctx context.Context, request openapi.GetVersionRequestObject) (openapi.GetVersionResponseObject, error) {
	return openapi.GetVersion200TextResponse(config.Version), nil
}
