package endec

import (
	"time"
)

type Purpose string

const (
	PurposePasswordReset Purpose = "password_reset"
	PurposeOAuthState    Purpose = "oauth_state"
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
