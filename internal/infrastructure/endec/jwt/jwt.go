package jwt

import (
	"crypto/rand"
	"io"
	"time"

	"github.com/Southclaws/fault"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
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

func New(cfg config.Config) (endec.EncrypterDecrypter, error) {
	if len(cfg.JWTSecret) == 0 {
		return nil, nil
	}

	return &jwtEncrypterDecrypter{key: cfg.JWTSecret}, nil
}

func (e *jwtEncrypterDecrypter) Encrypt(purpose endec.Purpose, data endec.Claims, lifespan time.Duration) (string, error) {
	if len(e.key) == 0 {
		return "", fault.New("no JWT secret provided")
	}

	if purpose == "" {
		return "", fault.New("no token purpose provided")
	}

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
	t.Header["typ"] = string(purpose)

	s, err := t.SignedString(e.key)
	if err != nil {
		return "", fault.Wrap(err)
	}

	return s, nil
}

func (e *jwtEncrypterDecrypter) Decrypt(purpose endec.Purpose, message string) (endec.Claims, error) {
	t, err := jwt.Parse(message, e.keyfunc, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, fault.Wrap(err)
	}

	if !t.Valid {
		return nil, fault.New("token flagged as invalid but no error was reported")
	}

	if typ, ok := t.Header["typ"].(string); !ok || typ != string(purpose) {
		return nil, fault.Newf("token was not issued for %s", purpose)
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
