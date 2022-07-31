package bindings

import (
	"context"
	"fmt"

	"github.com/Southclaws/storyden/backend/internal/config"
	"github.com/Southclaws/storyden/backend/pkg/transports/http/openapi"
)

type Version struct{}

func NewVersion() Version { return Version{} }

func (v *Version) GetVersion(ctx context.Context, request openapi.GetVersionRequestObject) any {
	fmt.Println(ctx, config.Version)
	return openapi.GetVersion200TextResponse(config.Version)
}
