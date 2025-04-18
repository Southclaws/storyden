package token

import (
	"encoding/json"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/rs/xid"
)

var (
	ErrTokenExpired = fault.New("token expired")
	ErrTokenRevoked = fault.New("token revoked")
)

type Token struct{ xid.ID }

func FromString(b string) (Token, error) {
	id, err := xid.FromString(b)
	if err != nil {
		return Token{}, err
	}

	return Token{id}, nil
}

func Generate() Token {
	return Token{xid.New()}
}

func (t Token) Bytes() []byte {
	return t.ID.Bytes()
}

func (t Token) String() string {
	return t.ID.String()
}

type Session struct {
	Token     Token                   `json:"t"`
	AccountID account.AccountID       `json:"a"`
	ExpiresAt time.Time               `json:"e"`
	RevokedAt opt.Optional[time.Time] `json:"r"`
}

type Validated Session

func (s *Session) Validate() (*Validated, error) {
	if s.RevokedAt.Ok() {
		return nil, ErrTokenRevoked
	}

	if s.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	return (*Validated)(s), nil
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
