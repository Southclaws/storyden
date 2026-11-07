package webpush

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

const (
	recordSize     = 4096
	headerSize     = 16 + 4 + 1 + 65
	maxPayloadSize = recordSize - headerSize - 1 - 16
)

var ErrPayloadTooLarge = errors.New("push payload must not exceed 3993 bytes")

func encrypt(subscription Subscription, payload []byte) ([]byte, error) {
	if len(payload) > maxPayloadSize {
		return nil, ErrPayloadTooLarge
	}

	if subscription.P256DH == nil {
		return nil, errors.New("push subscription must include a public key")
	}

	if len(subscription.Auth) != 16 {
		return nil, errors.New("push subscription must include a 16-byte authentication secret")
	}

	private, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	var salt [16]byte
	if _, err := rand.Read(salt[:]); err != nil {
		return nil, err
	}

	return encryptRecord(subscription, payload, private, salt)
}

func encryptRecord(
	subscription Subscription,
	payload []byte,
	private *ecdh.PrivateKey,
	salt [16]byte,
) ([]byte, error) {
	shared, err := private.ECDH(subscription.P256DH)
	if err != nil {
		return nil, err
	}

	public := private.PublicKey().Bytes()
	info := "WebPush: info\x00" + string(subscription.P256DH.Bytes()) + string(public)
	material, err := hkdf.Key(sha256.New, shared, subscription.Auth, info, 32)
	if err != nil {
		return nil, err
	}

	key, err := hkdf.Key(sha256.New, material, salt[:], "Content-Encoding: aes128gcm\x00", 16)
	if err != nil {
		return nil, err
	}

	nonce, err := hkdf.Key(sha256.New, material, salt[:], "Content-Encoding: nonce\x00", 12)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	header := make([]byte, 0, headerSize+len(payload)+1+gcm.Overhead())
	header = append(header, salt[:]...)
	header = binary.BigEndian.AppendUint32(header, recordSize)
	header = append(header, byte(len(public)))
	header = append(header, public...)

	plaintext := make([]byte, len(payload)+1)
	copy(plaintext, payload)
	plaintext[len(payload)] = 0x02

	return gcm.Seal(header, nonce, plaintext, nil), nil
}
