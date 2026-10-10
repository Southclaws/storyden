package bindings

import (
	"encoding/json"
	"testing"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/stretchr/testify/require"
)

func TestSerialiseWebAuthnCredentialCreationExtensions(t *testing.T) {
	var credential protocol.CredentialCreation
	credential.Response.Extensions = protocol.AuthenticationExtensions{
		CredProps: true,
		PRF:       &protocol.PRFInputs{},
		CredBlob:  protocol.URLEncodedBase64{0xfb, 0xff},
		Extra:     map[string]any{"customExtension": "value"},
	}

	options, err := serialiseWebAuthnCredentialCreationOptions(credential)
	require.NoError(t, err)
	encoded, err := json.Marshal(options.PublicKey.Extensions)
	require.NoError(t, err)
	require.JSONEq(t, `{"credProps":true,"prf":{},"credBlob":"-_8","customExtension":"value"}`, string(encoded))
}

func TestSerialiseWebAuthnCredentialCreationRejectsConflictingExtensions(t *testing.T) {
	var credential protocol.CredentialCreation
	credential.Response.Extensions = protocol.AuthenticationExtensions{
		CredProps: true,
		Extra:     map[string]any{"credProps": false},
	}

	_, err := serialiseWebAuthnCredentialCreationOptions(credential)
	require.Error(t, err)
}
