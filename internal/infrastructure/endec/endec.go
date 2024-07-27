package endec

import (
	"time"
)

type Encrypter interface {
	Encrypt(data Claims, lifespan time.Duration) (string, error)
}

type Decrypter interface {
	Decrypt(message string) (Claims, error)
}

type EncrypterDecrypter interface {
	Encrypter
	Decrypter
}

type Claims map[string]any
