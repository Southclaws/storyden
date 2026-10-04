package oauth

import (
	"context"
	"net/url"
	"strings"

	"github.com/Southclaws/storyden/app/resources/settings"
)

var clientAssertionSigningAlgorithms = []string{"RS256", "RS384", "RS512", "PS256", "PS384", "PS512", "ES256", "ES384", "ES512", "EdDSA"}

type Discovery struct {
	TokenEndpointAuthSigningAlgValuesSupported []string
	Issuer                                     string
	AuthorizationEndpoint                      string
	DeviceAuthorizationEndpoint                string
	TokenEndpoint                              string
	UserinfoEndpoint                           string
	RegistrationEndpoint                       string
	RegistrationModesSupported                 []string
	JWKSURI                                    string
	ResponseTypesSupported                     []string
	GrantTypesSupported                        []string
	CodeChallengeMethodsSupported              []string
	ScopesSupported                            []string
	SubjectTypesSupported                      []string
	IDTokenSigningAlgValuesSupported           []string
	TokenEndpointAuthMethodsSupported          []string
	ClientIDMetadataDocumentSupported          bool
}

func (s *Service) Discovery(ctx context.Context) (Discovery, error) {
	configuration, err := s.registrationSettings(ctx)
	if err != nil {
		return Discovery{}, err
	}

	endpointBase := s.apiEndpointBase()

	registrationEndpoint := ""
	if configuration.DynamicRegistrationEnabled.OrZero() {
		registrationEndpoint = endpointBase.JoinPath("oauth", "register").String()
	}

	var registrationModes []string
	if configuration.DynamicRegistrationEnabled.OrZero() {
		registrationModes = []string{"immediate"}
		if configuration.AutonomousRegistrationMode.OrZero() == settings.OAuthAutonomousRegistrationModeApproval {
			registrationModes = append(registrationModes, "approval")
		}
	}

	return Discovery{
		TokenEndpointAuthSigningAlgValuesSupported: clientAssertionSigningAlgorithms,
		Issuer:                            s.issuer,
		AuthorizationEndpoint:             endpointBase.JoinPath("oauth", "authorize").String(),
		DeviceAuthorizationEndpoint:       endpointBase.JoinPath("oauth", "device_authorization").String(),
		TokenEndpoint:                     endpointBase.JoinPath("oauth", "token").String(),
		UserinfoEndpoint:                  endpointBase.JoinPath("oauth", "userinfo").String(),
		RegistrationEndpoint:              registrationEndpoint,
		RegistrationModesSupported:        registrationModes,
		JWKSURI:                           endpointBase.JoinPath("oauth", "jwks").String(),
		ResponseTypesSupported:            []string{"code"},
		GrantTypesSupported:               []string{GrantTypeAuthorizationCode, GrantTypeRefreshToken, GrantTypeClientCredentials, GrantTypeDeviceCode},
		CodeChallengeMethodsSupported:     []string{CodeChallengeMethodS256},
		ScopesSupported:                   supportedScopes(),
		SubjectTypesSupported:             []string{"public"},
		IDTokenSigningAlgValuesSupported:  []string{"RS256"},
		TokenEndpointAuthMethodsSupported: []string{TokenEndpointAuthMethodNone, TokenEndpointAuthMethodClientSecretBasic, TokenEndpointAuthMethodClientSecretPost, TokenEndpointAuthMethodPrivateKeyJWT},
		ClientIDMetadataDocumentSupported: s.cimdEnabled(),
	}, nil
}

func (s *Service) apiEndpointBase() *url.URL {
	u := s.cfg.PublicAPIAddress
	path := strings.TrimRight(u.Path, "/")
	if !strings.HasSuffix(path, "/api") {
		return u.JoinPath("api")
	}

	return &u
}

// Issuer returns the OAuth issuer identifier (the authorization server base URL).
// This is derived from PublicAPIAddress with any trailing /api suffix removed.
func (s *Service) Issuer() string {
	return s.issuer
}

type JWK struct {
	Kty string
	Use string
	Alg string
	Kid string
	N   string
	E   string
}

func (s *Service) JWKS() []JWK {
	if s.signer == nil {
		return nil
	}

	return []JWK{{
		Kty: "RSA",
		Use: "sig",
		Alg: "RS256",
		Kid: s.kid,
		N:   b64url(s.signer.PublicKey.N.Bytes()),
		E:   b64url(bigEndian(s.signer.PublicKey.E)),
	}}
}

func bigEndian(v int) []byte {
	out := []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
	for len(out) > 1 && out[0] == 0 {
		out = out[1:]
	}

	return out
}
