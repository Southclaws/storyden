package access_key_test

import (
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/authentication/access_key"
)

func TestParseAccessKey(t *testing.T) {
	t.Parallel()

	t.Run("valid_personal_access_key", func(t *testing.T) {
		t.Parallel()
		key := "sdpak_010203040102030405060708091011121314151612344321"
		record, err := access_key.ParseAccessKeyToken(key)
		require.NoError(t, err)
		assert.Equal(t, access_key.AccessKeyKindPersonal, record.GetKind())
		assert.Equal(t, access_key.AccessKeyID("010203040102"), record.GetID())
		assert.Equal(t, access_key.AccessKeySecret("030405060708091011121314151612344321"), record.GetSecret())
	})

	t.Run("valid_bot_access_key", func(t *testing.T) {
		t.Parallel()
		key := "sdbak_010203040102030405060708091011121314151612344321"
		record, err := access_key.ParseAccessKeyToken(key)
		require.NoError(t, err)
		assert.Equal(t, access_key.AccessKeyKindBot, record.GetKind())
		assert.Equal(t, access_key.AccessKeyID("010203040102"), record.GetID())
		assert.Equal(t, access_key.AccessKeySecret("030405060708091011121314151612344321"), record.GetSecret())
	})

	t.Run("invalid_length", func(t *testing.T) {
		t.Parallel()
		key := "sdpak_short"
		_, err := access_key.ParseAccessKeyToken(key)
		assert.Error(t, err)
	})

	t.Run("invalid_kind", func(t *testing.T) {
		t.Parallel()
		key := "sdxxx_01020304010201020304050607080910111213141516"
		_, err := access_key.ParseAccessKeyToken(key)
		assert.Error(t, err)
	})

	t.Run("invalid_identifier_format", func(t *testing.T) {
		t.Parallel()
		key := "sdpak_010203;4010201020304050607080910111213141516"
		_, err := access_key.ParseAccessKeyToken(key)
		assert.Error(t, err)
	})

	t.Run("invalid_secret_format", func(t *testing.T) {
		t.Parallel()
		key := "sdpak_010203040102zzzzzzzzzz;zzzzzzzzzzzzzzzzzzzzz"
		_, err := access_key.ParseAccessKeyToken(key)
		assert.Error(t, err)
	})
}

func TestAccessKeyToken_Validate(t *testing.T) {
	t.Parallel()

	t.Run("valid_secret", func(t *testing.T) {
		t.Parallel()

		// Create a personal access key without expiry
		pak := access_key.NewPersonalAccessKey(opt.NewEmpty[time.Time]())

		// Get the string representation (what user would put in Authorization header)
		keyString := pak.String()

		// Parse the token from the string (simulates incoming request)
		token, err := access_key.ParseAccessKeyToken(keyString)
		require.NoError(t, err)

		// Validate the token against the stored record
		validatedToken, err := token.Validate(pak.AccessKeyRecord)
		require.NoError(t, err)
		assert.NotNil(t, validatedToken)
		assert.Equal(t, pak.KeyID, validatedToken.GetID())
		assert.Equal(t, pak.Kind, validatedToken.GetKind())
	})

	t.Run("invalid_secret", func(t *testing.T) {
		t.Parallel()

		// Create a personal access key
		pak := access_key.NewPersonalAccessKey(opt.NewEmpty[time.Time]())

		// Create a malformed key string with wrong secret
		wrongSecret := access_key.NewAccessKeySecret()
		wrongKeyString := pak.Kind.String() + "_" + string(pak.KeyID) + string(wrongSecret)

		// Parse the token from the wrong string
		token, err := access_key.ParseAccessKeyToken(wrongKeyString)
		require.NoError(t, err)

		// Validation should fail due to secret mismatch
		validatedToken, err := token.Validate(pak.AccessKeyRecord)
		assert.Error(t, err)
		assert.Nil(t, validatedToken)
		assert.Contains(t, err.Error(), "access key malformed")
	})

	t.Run("expired_key", func(t *testing.T) {
		t.Parallel()

		// Create a personal access key that expired 1 hour ago
		expiredTime := time.Now().Add(-1 * time.Hour)
		expiredPak := access_key.NewPersonalAccessKey(opt.New(expiredTime))

		// Get the string representation
		keyString := expiredPak.String()

		// Parse the token from the string
		token, err := access_key.ParseAccessKeyToken(keyString)
		require.NoError(t, err)

		// Validation should fail due to expiry
		validatedToken, err := token.Validate(expiredPak.AccessKeyRecord)
		assert.Error(t, err)
		assert.Nil(t, validatedToken)
		assert.Contains(t, err.Error(), "access key has expired")
	})

	t.Run("not_expired_key", func(t *testing.T) {
		t.Parallel()

		// Create a personal access key that expires in 1 hour
		futureTime := time.Now().Add(1 * time.Hour)
		futurePak := access_key.NewPersonalAccessKey(opt.New(futureTime))

		// Get the string representation
		keyString := futurePak.String()

		// Parse the token from the string
		token, err := access_key.ParseAccessKeyToken(keyString)
		require.NoError(t, err)

		// Validation should succeed
		validatedToken, err := token.Validate(futurePak.AccessKeyRecord)
		require.NoError(t, err)
		assert.NotNil(t, validatedToken)
		assert.Equal(t, futurePak.KeyID, validatedToken.GetID())
	})

	t.Run("bot_key_valid_secret", func(t *testing.T) {
		t.Parallel()

		// Create a bot access key without expiry
		bak := access_key.NewBotAccessKey(opt.NewEmpty[time.Time]())

		// Get the string representation
		keyString := bak.String()

		// Parse the token from the string
		token, err := access_key.ParseAccessKeyToken(keyString)
		require.NoError(t, err)

		// Validate the token against the stored record
		validatedToken, err := token.Validate(bak.AccessKeyRecord)
		require.NoError(t, err)
		assert.NotNil(t, validatedToken)
		assert.Equal(t, bak.KeyID, validatedToken.GetID())
		assert.Equal(t, access_key.AccessKeyKindBot, validatedToken.GetKind())
	})

	t.Run("bot_key_invalid_secret", func(t *testing.T) {
		t.Parallel()

		// Create a bot access key
		bak := access_key.NewBotAccessKey(opt.NewEmpty[time.Time]())

		// Create a malformed key string with wrong secret
		wrongSecret := access_key.NewAccessKeySecret()
		wrongKeyString := bak.Kind.String() + "_" + string(bak.KeyID) + string(wrongSecret)

		// Parse the token from the wrong string
		token, err := access_key.ParseAccessKeyToken(wrongKeyString)
		require.NoError(t, err)

		// Validation should fail due to secret mismatch
		validatedToken, err := token.Validate(bak.AccessKeyRecord)
		assert.Error(t, err)
		assert.Nil(t, validatedToken)
		assert.Contains(t, err.Error(), "access key malformed")
	})

	t.Run("no_expiry_key", func(t *testing.T) {
		t.Parallel()

		// Create a personal access key with no expiry (explicitly test empty optional)
		noExpiryPak := access_key.NewPersonalAccessKey(opt.NewEmpty[time.Time]())

		// Get the string representation
		keyString := noExpiryPak.String()

		// Parse the token from the string
		token, err := access_key.ParseAccessKeyToken(keyString)
		require.NoError(t, err)

		// Validation should succeed since no expiry means never expires
		validatedToken, err := token.Validate(noExpiryPak.AccessKeyRecord)
		require.NoError(t, err)
		assert.NotNil(t, validatedToken)
		assert.Equal(t, noExpiryPak.KeyID, validatedToken.GetID())
	})

	t.Run("far_future_expiry", func(t *testing.T) {
		t.Parallel()

		// Create a personal access key that expires far in the future
		farFutureTime := time.Now().Add(365 * 24 * time.Hour) // 1 year from now
		farFuturePak := access_key.NewPersonalAccessKey(opt.New(farFutureTime))

		// Get the string representation
		keyString := farFuturePak.String()

		// Parse the token from the string
		token, err := access_key.ParseAccessKeyToken(keyString)
		require.NoError(t, err)

		// Validation should succeed
		validatedToken, err := token.Validate(farFuturePak.AccessKeyRecord)
		require.NoError(t, err)
		assert.NotNil(t, validatedToken)
		assert.Equal(t, farFuturePak.KeyID, validatedToken.GetID())
	})

	t.Run("edge_case_just_expired", func(t *testing.T) {
		t.Parallel()

		// Create a personal access key that expires very recently (within the last second)
		justExpiredTime := time.Now().Add(-1 * time.Millisecond)
		justExpiredPak := access_key.NewPersonalAccessKey(opt.New(justExpiredTime))

		// Get the string representation
		keyString := justExpiredPak.String()

		// Parse the token from the string
		token, err := access_key.ParseAccessKeyToken(keyString)
		require.NoError(t, err)

		// Validation should fail due to expiry
		validatedToken, err := token.Validate(justExpiredPak.AccessKeyRecord)
		assert.Error(t, err)
		assert.Nil(t, validatedToken)
		assert.Contains(t, err.Error(), "access key has expired")
	})
}

func TestAccessKeyRecordFromAuthenticationRecord(t *testing.T) {
	t.Parallel()

	t.Run("valid_record", func(t *testing.T) {
		t.Parallel()
		pak := access_key.NewPersonalAccessKey(opt.NewEmpty[time.Time]())

		auth := authentication.Authentication{
			Identifier: pak.GetAuthenticationRecordIdentifier(),
			Token:      string(pak.Hash),
		}
		record, err := access_key.AccessKeyRecordFromAuthenticationRecord(auth)
		require.NoError(t, err)
		assert.Equal(t, access_key.AccessKeyKindPersonal, record.Kind)
		assert.Equal(t, access_key.AccessKeyID(pak.KeyID), record.KeyID)

		header := pak.String()
		akt, err := access_key.ParseAccessKeyToken(header)
		require.NoError(t, err)

		vakt, err := akt.Validate(*record)
		require.NoError(t, err)
		assert.Equal(t, pak.KeyID, vakt.AccessKeyToken.GetID())
	})

	t.Run("mismatched_identifier", func(t *testing.T) {
		t.Parallel()
		pak := access_key.NewPersonalAccessKey(opt.NewEmpty[time.Time]())

		auth := authentication.Authentication{
			Identifier: string("wrong"),
			Token:      string(pak.Hash),
		}

		_, err := access_key.AccessKeyRecordFromAuthenticationRecord(auth)
		assert.Error(t, err)
	})
}
