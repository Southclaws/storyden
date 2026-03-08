package plugin_auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math/big"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/rs/xid"
	"golang.org/x/crypto/nacl/secretbox"

	"github.com/Southclaws/storyden/app/resources/plugin"
)

var ErrInvalidToken = errors.New("invalid token")

const (
	secretLength        = 32
	nonceLength         = 24
	charset             = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	ExternalTokenPrefix = "sdprt_"
)

func GenerateSecret() (string, error) {
	secret := make([]byte, secretLength)
	max := big.NewInt(int64(len(charset)))

	for i := range secret {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fault.Wrap(err, fmsg.With("failed to generate random number"))
		}
		secret[i] = charset[num.Int64()]
	}

	return string(secret), nil
}

func GenerateExternalToken() (string, error) {
	secret, err := GenerateSecret()
	if err != nil {
		return "", err
	}

	return ExternalTokenPrefix + secret, nil
}

func IsExternalToken(token string) bool {
	if !strings.HasPrefix(token, ExternalTokenPrefix) {
		return false
	}

	raw := strings.TrimPrefix(token, ExternalTokenPrefix)
	if len(raw) != secretLength {
		return false
	}

	for _, c := range raw {
		if !strings.ContainsRune(charset, c) {
			return false
		}
	}

	return true
}

func SealToken(pluginID plugin.InstallationID, secret string) (string, error) {
	if len(secret) != secretLength {
		return "", fault.Newf("secret must be %d bytes", secretLength)
	}

	var secretKey [32]byte
	copy(secretKey[:], []byte(secret))

	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", fault.Wrap(err, fmsg.With("failed to generate nonce"))
	}

	message := []byte(pluginID.String())
	encrypted := secretbox.Seal(nonce[:], message, &nonce, &secretKey)

	return base64.RawURLEncoding.EncodeToString(encrypted), nil
}

func OpenToken(token string, secret string) (plugin.InstallationID, error) {
	if len(secret) != secretLength {
		return plugin.InstallationID(xid.NilID()), fault.Newf("secret must be %d bytes", secretLength)
	}

	var secretKey [32]byte
	copy(secretKey[:], []byte(secret))

	encrypted, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return plugin.InstallationID(xid.NilID()), fault.Wrap(err, fmsg.With("failed to decode token"))
	}

	if len(encrypted) < nonceLength {
		return plugin.InstallationID(xid.NilID()), fault.New("token too short")
	}

	var nonce [24]byte
	copy(nonce[:], encrypted[:nonceLength])

	decrypted, ok := secretbox.Open(nil, encrypted[nonceLength:], &nonce, &secretKey)
	if !ok {
		return plugin.InstallationID(xid.NilID()), fault.New("failed to decrypt token")
	}

	formatted, err := xid.FromString(string(decrypted))
	if err != nil {
		return plugin.InstallationID(xid.NilID()), fault.New("failed to parse id")
	}

	return plugin.InstallationID(formatted), nil
}
