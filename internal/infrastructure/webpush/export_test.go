package webpush

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
)

func NewTestClient(t *testing.T, transport http.RoundTripper) (config.Config, *Client) {
	t.Helper()

	cfg := pushConfig(t)
	client, err := New(cfg)
	require.NoError(t, err)

	client.http.Transport = transport

	return cfg, client
}

func NewTestSubscription(t *testing.T, endpoint string) (SubscriptionInput, *ecdh.PrivateKey) {
	t.Helper()

	private, err := ecdh.P256().GenerateKey(rand.Reader)
	require.NoError(t, err)

	return SubscriptionInput{
		Endpoint: endpoint,
		P256DH:   base64.RawURLEncoding.EncodeToString(private.PublicKey().Bytes()),
		Auth:     "AAAAAAAAAAAAAAAAAAAAAA",
	}, private
}

func DecryptTestPayload(
	t *testing.T,
	private *ecdh.PrivateKey,
	auth []byte,
	body []byte,
) ([]byte, error) {
	t.Helper()

	r := require.New(t)
	r.GreaterOrEqual(len(body), 103)
	r.Equal(uint32(4096), binary.BigEndian.Uint32(body[16:20]))
	r.Equal(byte(65), body[20])

	public, err := ecdh.P256().NewPublicKey(body[21:86])
	r.NoError(err)

	shared, err := private.ECDH(public)
	r.NoError(err)

	info := "WebPush: info\x00" + string(private.PublicKey().Bytes()) + string(public.Bytes())
	material, err := hkdf.Key(sha256.New, shared, auth, info, 32)
	r.NoError(err)

	key, err := hkdf.Key(sha256.New, material, body[:16], "Content-Encoding: aes128gcm\x00", 16)
	r.NoError(err)

	nonce, err := hkdf.Key(sha256.New, material, body[:16], "Content-Encoding: nonce\x00", 12)
	r.NoError(err)

	block, err := aes.NewCipher(key)
	r.NoError(err)

	gcm, err := cipher.NewGCM(block)
	r.NoError(err)

	return gcm.Open(nil, nonce, body[86:], nil)
}
