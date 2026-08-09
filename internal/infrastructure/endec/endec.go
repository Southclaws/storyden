package endec

import (
	"time"
)

// Purpose is the RFC 7519 "typ" header value, explicit typing per RFC 8725
// section 3.11 is what stops one feature's token verifying as another's
type Purpose string

const (
	PurposePasswordReset Purpose = "storyden-password-reset+jwt"
	PurposeOAuthState    Purpose = "storyden-oauth-state+jwt"
)

type Encrypter interface {
	Encrypt(purpose Purpose, data Claims, lifespan time.Duration) (string, error)
}

type Decrypter interface {
	Decrypt(purpose Purpose, message string) (Claims, error)
}

type EncrypterDecrypter interface {
	Encrypter
	Decrypter
}

type Claims map[string]any
