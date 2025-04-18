package jwt

import (
	"crypto/rand"
	"io"
	"time"

	"github.com/Southclaws/fault"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/infrastructure/endec"
)

type jwtEncrypterDecrypter struct {
	key []byte
}

func Build() fx.Option {
	return fx.Provide(
		fx.Annotate(New,
			fx.As(new(endec.EncrypterDecrypter)),
			fx.As(new(endec.Encrypter)),
			fx.As(new(endec.Decrypter)),
		),
	)
}

func New() (endec.EncrypterDecrypter, error) {
	// NOTE: This is currently not compatible when running replicas of the pod
	// in a cluster, as the key is generated randomly each boot.
	// TODO: Move this key generation to Settings init and load key from there.
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fault.Wrap(err)
	}

	return &jwtEncrypterDecrypter{key: key}, nil
}

func (e *jwtEncrypterDecrypter) Encrypt(data endec.Claims, lifespan time.Duration) (string, error) {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", fault.Wrap(err)
	}

	expires := time.Now().UTC().Add(lifespan)

	claims := jwt.MapClaims{
		"jti": nonce,
		"exp": jwt.NewNumericDate(expires),
	}

	for k, v := range data {
		claims[k] = v
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := t.SignedString(e.key)
	if err != nil {
		return "", fault.Wrap(err)
	}

	return s, nil
}

func (e *jwtEncrypterDecrypter) Decrypt(message string) (endec.Claims, error) {
	t, err := jwt.Parse(message, e.keyfunc)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	if !t.Valid {
		return nil, fault.New("token flagged as invalid but no error was reported")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fault.New("invalid token")
	}

	return endec.Claims(claims), nil
}

func (e *jwtEncrypterDecrypter) keyfunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fault.Newf("invalid jwt algorithm %e", token.Header["alg"])
	}

	return e.key, nil
}
