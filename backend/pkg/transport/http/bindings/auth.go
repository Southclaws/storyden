package bindings

import (
	"net/http"

	"github.com/Southclaws/storyden/backend/pkg/services/authentication"
	"github.com/Southclaws/storyden/backend/pkg/transport/http/openapi"
)

type Authentication struct{ Service authentication.Service }

func NewAuthentication(s authentication.Service) Authentication { return Authentication{s} }

func (i *Authentication) GetV1AuthPassword(w http.ResponseWriter, r *http.Request, params openapi.GetV1AuthPasswordParams) {
	//
}
