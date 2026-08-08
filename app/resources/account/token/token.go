package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
)

var (
	ErrTokenExpired = fault.New("token expired")
	ErrTokenRevoked = fault.New("token revoked")
	ErrTokenInvalid = fault.New("token malformed", ftag.With(ftag.Unauthenticated))
)

const (
	secretBytes  = 32
	secretLength = 43
)

// Token is the bearer secret held by the client. It is never persisted, only
// its hash is, so a database or cache leak cannot be replayed as a session.
type Token struct{ secret string }

func Generate() (Token, error) {
	b := make([]byte, secretBytes)
	if _, err := rand.Read(b); err != nil {
		return Token{}, fault.Wrap(err)
	}

	return Token{base64.RawURLEncoding.EncodeToString(b)}, nil
}

func FromString(s string) (Token, error) {
	if len(s) != secretLength {
		return Token{}, fault.Wrap(ErrTokenInvalid)
	}

	if _, err := base64.RawURLEncoding.DecodeString(s); err != nil {
		return Token{}, fault.Wrap(ErrTokenInvalid)
	}

	return Token{s}, nil
}

func (t Token) String() string {
	return t.secret
}

func (t Token) Hash() string {
	sum := sha256.Sum256([]byte(t.secret))
	return hex.EncodeToString(sum[:])
}

type Session struct {
	TokenHash string                  `json:"h"`
	AccountID account.AccountID       `json:"a"`
	ExpiresAt time.Time               `json:"e"`
	RevokedAt opt.Optional[time.Time] `json:"r"`
}

// Issued pairs a new session with the one and only copy of its secret.
type Issued struct {
	Token   Token
	Session Session
}

type Validated Session

func (s Session) Validate() (*Validated, error) {
	if s.RevokedAt.Ok() {
		return nil, ErrTokenRevoked
	}

	if s.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	return (*Validated)(&s), nil
}

func (t Session) Serialise() ([]byte, error) {
	jsonData, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func Deserialise(data []byte) (*Session, error) {
	t := Session{}
	err := json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
