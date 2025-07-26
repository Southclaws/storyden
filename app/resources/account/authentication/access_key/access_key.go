package access_key

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/alexedwards/argon2id"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
)

var (
	errInvalidAccessKey               = fault.New("invalid access key")
	errMalformedStoredAccessKeyRecord = fault.New("malformed stored access key record")
	errRejectedAccessKey              = fault.New("access key rejected")
	errRevoked                        = fault.New("access key revoked")
)

const (
	kindLength       = 5                                              // sdpak, sdbak
	prefixLength     = 6                                              // sdpak_ or sdbak_
	identifierLength = 12                                             // random identifier
	secretLength     = 32                                             // random secret
	AccessKeyLength  = prefixLength + identifierLength + secretLength // 46 characters
)

//go:generate go run github.com/Southclaws/enumerator

type accessKeyKindEnum string

const (
	accessKeyKindPersonal accessKeyKindEnum = `sdpak`
	accessKeyKindBot      accessKeyKindEnum = `sdbak`
)

type AccessKeyID string

func NewAccessKeyID() AccessKeyID {
	return AccessKeyID(randomString(identifierLength))
}

type AccessKeySecret string

func NewAccessKeySecret() AccessKeySecret {
	return AccessKeySecret(randomString(secretLength))
}

// Represents an input token from a request, parsed but not validated.
type AccessKeyToken struct {
	kind   AccessKeyKind
	id     AccessKeyID
	secret AccessKeySecret
}

func (k *AccessKeyToken) GetAuthenticationRecordIdentifier() string {
	return fmt.Sprintf("%s_%s", k.kind, k.id)
}

func (t *AccessKeyToken) GetKind() AccessKeyKind {
	return t.kind
}

func (t *AccessKeyToken) GetID() AccessKeyID {
	return t.id
}

func (t *AccessKeyToken) GetSecret() AccessKeySecret {
	return t.secret
}

type ValidatedAccessKeyToken struct {
	AccessKeyToken
}

type AccessKeyHash string

// Represents the actual stored key with its hash, no secret.
type AccessKeyRecord struct {
	Kind      AccessKeyKind
	AuthID    authentication.ID
	KeyID     AccessKeyID
	Hash      AccessKeyHash
	CreatedAt time.Time
	Expires   opt.Optional[time.Time]
	Disabled  bool
}

func (r *AccessKeyRecord) GetAuthenticationRecordIdentifier() string {
	return fmt.Sprintf("%s_%s", r.Kind, r.KeyID)
}

// Expose this to the member who created the access key once, do not store.
type AccessKeyRecordWithSecret struct {
	AccessKeyRecord
	Name   string
	secret AccessKeySecret
}

func (a AccessKeyRecordWithSecret) String() string {
	return fmt.Sprintf("%s_%s%s", a.Kind, a.KeyID, a.secret)
}

func newAccessKey(kind AccessKeyKind, expiry opt.Optional[time.Time]) AccessKeyRecordWithSecret {
	secret := NewAccessKeySecret()
	hash, err := argon2id.CreateHash(string(secret), argon2id.DefaultParams)
	if err != nil {
		panic(err)
	}

	return AccessKeyRecordWithSecret{
		AccessKeyRecord: AccessKeyRecord{
			Kind:      kind,
			KeyID:     NewAccessKeyID(),
			Hash:      AccessKeyHash(hash),
			CreatedAt: time.Now(),
			Expires:   expiry,
		},
		secret: secret,
	}
}

func NewPersonalAccessKey(expiry opt.Optional[time.Time]) AccessKeyRecordWithSecret {
	return newAccessKey(AccessKeyKindPersonal, expiry)
}

func NewBotAccessKey(expiry opt.Optional[time.Time]) AccessKeyRecordWithSecret {
	return newAccessKey(AccessKeyKindBot, expiry)
}

// AccessKeyIdentifier represents the first 2 parts of an access key:
// "sdpak_12345678afbe"
type AccessKeyIdentifier struct {
	kind AccessKeyKind
	id   AccessKeyID
}

func ParseAccessKeyIdentifier(s string) (*AccessKeyIdentifier, error) {
	if len(s) < prefixLength+identifierLength {
		return nil, fault.Wrap(errInvalidAccessKey, fmsg.With("incorrect identifier length"), ftag.With(ftag.InvalidArgument))
	}

	rawkind := s[:kindLength]
	rawid := s[prefixLength : prefixLength+identifierLength]

	kind, err := NewAccessKeyKind(rawkind)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("invalid kind"), ftag.With(ftag.InvalidArgument))
	}

	if !isValidRandomString(rawid) {
		return nil, fault.Wrap(errInvalidAccessKey, fmsg.With("invalid identifier format"), ftag.With(ftag.InvalidArgument))
	}

	return &AccessKeyIdentifier{
		kind: kind,
		id:   AccessKeyID(rawid),
	}, nil
}

func ParseAccessKeyToken(s string) (*AccessKeyToken, error) {
	if len(s) < prefixLength+identifierLength+secretLength {
		return nil, fault.Wrap(errInvalidAccessKey, fmsg.With("incorrect token length"), ftag.With(ftag.InvalidArgument))
	}

	aki, err := ParseAccessKeyIdentifier(s)
	if err != nil {
		return nil, err
	}

	rawsecret := s[prefixLength+identifierLength:]

	if !isValidRandomString(rawsecret) {
		return nil, fault.Wrap(errInvalidAccessKey, fmsg.With("invalid secret format"), ftag.With(ftag.InvalidArgument))
	}

	return &AccessKeyToken{
		kind:   aki.kind,
		id:     aki.id,
		secret: AccessKeySecret(rawsecret),
	}, nil
}

func (k *AccessKeyToken) Validate(r AccessKeyRecord) (*ValidatedAccessKeyToken, error) {
	// Perform the hash even if disabled or expired to prevent timing attacks.
	ok, _, err := argon2id.CheckHash(string(k.secret), string(r.Hash))
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fault.Wrap(errInvalidAccessKey,
			fmsg.WithDesc("access key malformed", "The provided access key is malformed or revoked."),
			ftag.With(ftag.PermissionDenied))
	}

	if expires, ok := r.Expires.Get(); ok {
		if expires.Before(time.Now()) {
			return nil, fault.Wrap(errRejectedAccessKey,
				fmsg.WithDesc("access key has expired", "The provided access key has expired."),
				ftag.With(ftag.PermissionDenied))
		}
	}

	if r.Disabled {
		return nil, fault.Wrap(errRevoked, ftag.With(ftag.PermissionDenied))
	}

	return &ValidatedAccessKeyToken{
		AccessKeyToken: *k,
	}, nil
}

func AccessKeyRecordFromAuthenticationRecord(a authentication.Authentication) (*AccessKeyRecord, error) {
	aki, err := ParseAccessKeyIdentifier(a.Identifier)
	if err != nil {
		return nil, err
	}

	hash := a.Token

	return &AccessKeyRecord{
		Kind:      aki.kind,
		AuthID:    a.ID,
		KeyID:     AccessKeyID(aki.id),
		Hash:      AccessKeyHash(hash),
		CreatedAt: a.Created,
		Expires:   a.Expires,
		Disabled:  a.Disabled,
	}, nil
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func randomString(n int) string {
	key := make([]byte, n)
	max := big.NewInt(int64(len(charset)))

	for i := range key {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		key[i] = charset[num.Int64()]
	}

	return string(key)
}

func isValidRandomString(s string) bool {
	for _, c := range s {
		for i := 0; i < len(charset); i++ {
			if c == rune(charset[i]) {
				break
			}
			if i == len(charset)-1 {
				return false
			}
		}
	}

	return true
}
