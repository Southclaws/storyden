package oauth

import "strings"

type Discovery struct {
	Issuer                           string
	AuthorizationEndpoint            string
	DeviceAuthorizationEndpoint      string
	TokenEndpoint                    string
	UserinfoEndpoint                 string
	JWKSURI                          string
	ResponseTypesSupported           []string
	GrantTypesSupported              []string
	CodeChallengeMethodsSupported    []string
	ScopesSupported                  []string
	SubjectTypesSupported            []string
	IDTokenSigningAlgValuesSupported []string
}

func (s *Service) Discovery() Discovery {
	endpointBase := strings.TrimSuffix(s.apiEndpointBase(), "/")

	return Discovery{
		Issuer:                           s.issuer,
		AuthorizationEndpoint:            endpointBase + "/oauth/authorize",
		DeviceAuthorizationEndpoint:      endpointBase + "/oauth/device_authorization",
		TokenEndpoint:                    endpointBase + "/oauth/token",
		UserinfoEndpoint:                 endpointBase + "/oauth/userinfo",
		JWKSURI:                          endpointBase + "/oauth/jwks",
		ResponseTypesSupported:           []string{"code"},
		GrantTypesSupported:              []string{GrantTypeAuthorizationCode, GrantTypeRefreshToken, GrantTypeClientCredentials, GrantTypeDeviceCode},
		CodeChallengeMethodsSupported:    []string{CodeChallengeMethodS256},
		ScopesSupported:                  supportedScopes(),
		SubjectTypesSupported:            []string{"public"},
		IDTokenSigningAlgValuesSupported: []string{"RS256"},
	}
}

func (s *Service) apiEndpointBase() string {
	u := s.cfg.PublicAPIAddress
	path := strings.TrimRight(u.Path, "/")
	if !strings.HasSuffix(path, "/api") {
		path += "/api"
	}
	u.Path = path

	return u.String()
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
