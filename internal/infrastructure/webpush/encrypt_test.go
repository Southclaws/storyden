package webpush

import (
	"bytes"
	"crypto/ecdh"
	"encoding/base64"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptionRFC8291(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	subscription, err := ParseSubscription(SubscriptionInput{
		Endpoint: "https://push.example.net/push/JzLQ3raZJfFBR0aqvOMsLrt54w4rJUsV",
		P256DH: "BCVxsr7N_eNgVRqvHtD0zTZsEc6-VV-JvLexhqUzORcx" +
			"aOzi6-AYWXvTBHm4bjyPjs7Vd8pZGH6SRpkNtoIAiw4",
		Auth: "BTBZMqHH6r4Tts7J_aSIgg",
	})
	r.NoError(err)

	privateBytes, err := base64.RawURLEncoding.DecodeString("yfWPiYE-n46HLnH0KqZOF1fJJU3MYrct3AELtAQ-oRw")
	r.NoError(err)

	private, err := ecdh.P256().NewPrivateKey(privateBytes)
	r.NoError(err)

	salt, err := base64.RawURLEncoding.DecodeString("DGv6ra1nlYgDCS1FRnbzlw")
	r.NoError(err)

	payload := []byte("When I grow up, I want to be a watermelon")
	body, err := encryptRecord(*subscription, payload, private, [16]byte(salt))
	r.NoError(err)

	expected := "DGv6ra1nlYgDCS1FRnbzlwAAEABBBP4z9KsN6nGRTbVYI_c7VJSPQTBtkgcy27ml" +
		"mlMoZIIgDll6e3vCYLocInmYWAmS6TlzAC8wEqKK6PBru3jl7A_yl95bQpu6cVPT" +
		"pK4Mqgkf1CXztLVBSt2Ks3oZwbuwXPXLWyouBWLVWGNWQexSgSxsj_Qulcy4a-fN"
	a.Equal(expected, base64.RawURLEncoding.EncodeToString(body))
	a.Len(body, 144)
	a.Equal("When I grow up, I want to be a watermelon", string(payload))
}

func TestEncryptionPayloadSizes(t *testing.T) {
	r := require.New(t)

	input, receiver := NewTestSubscription(t, "https://push.example/message")
	subscription, err := ParseSubscription(input)
	r.NoError(err)

	for _, size := range []int{0, 1, 3993, 3994} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			r := require.New(t)
			a := assert.New(t)

			payload := bytes.Repeat([]byte{'a'}, size)
			body, err := encrypt(*subscription, payload)
			if size > 3993 {
				r.ErrorIs(err, ErrPayloadTooLarge)
				a.Empty(body)
				return
			}

			r.NoError(err)
			a.Len(body, size+103)

			plaintext, err := DecryptTestPayload(t, receiver, subscription.Auth, body)
			r.NoError(err)
			a.Equal(append(payload, 0x02), plaintext)
		})
	}
}

func TestEncryptionFreshKeysAndSalt(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	input, receiver := NewTestSubscription(t, "https://push.example/message")
	subscription, err := ParseSubscription(input)
	r.NoError(err)

	payload := []byte("notification-id")
	first, err := encrypt(*subscription, payload)
	r.NoError(err)

	second, err := encrypt(*subscription, payload)
	r.NoError(err)

	a.NotEqual(first[:16], second[:16])
	a.NotEqual(first[21:86], second[21:86])

	first[len(first)-1] ^= 1
	_, err = DecryptTestPayload(t, receiver, subscription.Auth, first)
	r.Error(err)
}
