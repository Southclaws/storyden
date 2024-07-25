package bindings

import (
	"context"

	"github.com/Southclaws/storyden/app/transports/openapi"
	"github.com/Southclaws/storyden/internal/saml"
)

type SAML struct {
	samlsp *saml.SAML
}

func NewSAML(sp *saml.SAML) SAML {
	return SAML{sp}
}

func (s *SAML) SAMLLogin(ctx context.Context, request openapi.SAMLLoginRequestObject) (openapi.SAMLLoginResponseObject, error) {
	return nil, nil
}

func (s *SAML) SAMLMetadataRead(ctx context.Context, request openapi.SAMLMetadataReadRequestObject) (openapi.SAMLMetadataReadResponseObject, error) {
	return nil, nil
}

func (s *SAML) SAMLACSGet(ctx context.Context, request openapi.SAMLACSGetRequestObject) (openapi.SAMLACSGetResponseObject, error) {
	return nil, nil
}
