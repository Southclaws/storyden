package securecookie

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"go.uber.org/fx"
	"golang.org/x/crypto/nacl/secretbox"

	"github.com/Southclaws/storyden/internal/config"
)

type Session struct {
	key [32]byte
}

func Build() fx.Option {
	return fx.Provide(New)
}

func New(cfg config.Config) (*Session, error) {
	sessionKey, err := hex.DecodeString(cfg.SessionKey)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to parse session key as a hexadecimal string"))
	}

	var secretKey [32]byte
	copy(secretKey[:], sessionKey)

	return &Session{
		key: secretKey,
	}, nil
}

func (s *Session) Encrypt(message string) string {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	b := secretbox.Seal(nonce[:], []byte(message), &nonce, &s.key)

	out := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(out, b)

	return string(out)
}

func (s *Session) Decrypt(message string) (string, bool) {
	box := make([]byte, hex.DecodedLen(len(message)))
	var decryptNonce [24]byte

	hex.Decode(box, []byte(message))
	if len(box) == 0 {
		return "", false
	}

	copy(decryptNonce[:], box[:24])

	result, ok := secretbox.Open(nil, box[24:], &decryptNonce, &s.key)
	if !ok {
		return "", false
	}

	return string(result), true
}
