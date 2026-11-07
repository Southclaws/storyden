package webpush

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/infrastructure/httpsafe"
)

func TestDeliveryRejectsPrivateAddresses(t *testing.T) {
	r := require.New(t)

	client, err := New(pushConfig(t))
	r.NoError(err)

	input, _ := NewTestSubscription(t, "https://push.example/message")
	subscription, err := ParseSubscription(input)
	r.NoError(err)

	for _, host := range []string{
		"127.0.0.1",
		"169.254.169.254",
		"10.0.0.1",
		"[::1]",
		"[::ffff:169.254.169.254]",
		"[fc00::1]",
		"[fec0::1]",
		"[64:ff9b:1::a00:1]",
		"[fc00::1%en0]",
	} {
		t.Run(host, func(t *testing.T) {
			subscription.Endpoint = url.URL{
				Scheme: "https",
				Host:   host,
				Path:   "/private",
			}

			err := client.Send(context.Background(), *subscription, []byte("notification-id"))
			require.ErrorIs(t, err, httpsafe.ErrDisallowedHost)
		})
	}
}

func TestDeliveryDoesNotFollowRedirects(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	client, err := New(pushConfig(t))
	r.NoError(err)

	input, _ := NewTestSubscription(t, "https://push.example/message")
	subscription, err := ParseSubscription(input)
	r.NoError(err)

	requests := 0
	client.http.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		requests++
		return &http.Response{
			StatusCode: http.StatusTemporaryRedirect,
			Header: http.Header{
				"Location": {"https://169.254.169.254/private"},
			},
			Body: http.NoBody,
		}, nil
	})

	err = client.Send(context.Background(), *subscription, []byte("notification-id"))
	r.ErrorContains(err, "HTTP 307")
	a.Equal(1, requests)
}

func TestDeliveryErrorsDoNotExposeEndpoint(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	client, err := New(pushConfig(t))
	r.NoError(err)

	input, _ := NewTestSubscription(t, "https://push.example/secret-capability?token=private-token")
	subscription, err := ParseSubscription(input)
	r.NoError(err)

	networkError := errors.New("connection failed")
	client.http.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, networkError
	})

	err = client.Send(context.Background(), *subscription, []byte("notification-id"))
	r.ErrorIs(err, networkError)
	a.NotContains(err.Error(), "secret-capability")
	a.NotContains(err.Error(), "private-token")
}

func TestEncryptionRejectsTampering(t *testing.T) {
	r := require.New(t)

	input, receiver := NewTestSubscription(t, "https://push.example/message")
	subscription, err := ParseSubscription(input)
	r.NoError(err)

	body, err := encrypt(*subscription, []byte("notification-id"))
	r.NoError(err)

	for _, tc := range []struct {
		name   string
		offset int
	}{
		{
			name:   "salt",
			offset: 0,
		},
		{
			name:   "ciphertext",
			offset: 86,
		},
		{
			name:   "authentication_tag",
			offset: len(body) - 1,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tampered := append([]byte(nil), body...)
			tampered[tc.offset] ^= 1

			_, err := DecryptTestPayload(t, receiver, subscription.Auth, tampered)
			require.Error(t, err)
		})
	}

	t.Run("wrong_authentication_secret", func(t *testing.T) {
		auth := append([]byte(nil), subscription.Auth...)
		auth[0] ^= 1

		_, err := DecryptTestPayload(t, receiver, auth, body)
		require.Error(t, err)
	})

	t.Run("wrong_receiver", func(t *testing.T) {
		_, other := NewTestSubscription(t, "https://other.example/message")

		_, err := DecryptTestPayload(t, other, subscription.Auth, body)
		require.Error(t, err)
	})
}

func FuzzParseSubscription(f *testing.F) {
	input := SubscriptionInput{
		Endpoint: "https://push.example/message",
		P256DH: "BCVxsr7N_eNgVRqvHtD0zTZsEc6-VV-JvLexhqUzORcx" +
			"aOzi6-AYWXvTBHm4bjyPjs7Vd8pZGH6SRpkNtoIAiw4",
		Auth: "BTBZMqHH6r4Tts7J_aSIgg",
	}
	f.Add(input.Endpoint, input.P256DH, input.Auth)
	f.Add("https://127.0.0.1/private", input.P256DH, input.Auth)
	f.Add("https://[::ffff:127.0.0.1]/private", input.P256DH, input.Auth)
	f.Add("https://[fc00::1%25en0]/private", input.P256DH, input.Auth)
	f.Add("https://push.example:999999/message", input.P256DH, input.Auth)
	f.Add("https://push.example/message", "AAAA", "AAAA")
	f.Add("https://push.example/message#fragment", input.P256DH, input.Auth)
	f.Add("https://push.example/message", strings.Repeat("A", 87), input.Auth)

	f.Fuzz(func(t *testing.T, endpoint, public, auth string) {
		subscription, err := ParseSubscription(SubscriptionInput{
			Endpoint: endpoint,
			P256DH:   public,
			Auth:     auth,
		})
		if err != nil {
			return
		}

		r := require.New(t)
		a := assert.New(t)

		r.NotNil(subscription.P256DH)
		a.Len(subscription.P256DH.Bytes(), 65)
		a.Len(subscription.Auth, 16)
		a.Equal("https", subscription.Endpoint.Scheme)
		a.NotEmpty(subscription.Endpoint.Hostname())
		a.Nil(subscription.Endpoint.User)
		a.Empty(subscription.Endpoint.Fragment)
		a.False(httpsafe.IsDisallowedHost(subscription.Endpoint.Hostname()))
	})
}
