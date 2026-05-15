package pubsub

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func TestInjectSessionContext_IncludesGuestPermissions(t *testing.T) {
	msg := message.NewMessage("test-uuid", []byte("test payload"))
	ctx := session.WithGuestPermissions(t.Context(), rbac.NewList(rbac.PermissionListProfiles))

	err := injectSessionContext(ctx, msg)

	require.NoError(t, err)
	assert.Empty(t, msg.Metadata.Get(accountIDKey))
	assert.JSONEq(t, `["LIST_PROFILES"]`, msg.Metadata.Get(permissionsKey))
}

func TestInjectSessionContext_IncludesPermissionsWithAccount(t *testing.T) {
	msg := message.NewMessage("test-uuid", []byte("test payload"))
	accountID := account.AccountID(xid.New())
	ctx := session.WithAccountPermissions(
		t.Context(),
		account.Account{ID: accountID},
		rbac.NewList(rbac.PermissionReadPublishedThreads),
	)

	err := injectSessionContext(ctx, msg)

	require.NoError(t, err)
	assert.Equal(t, accountID.String(), msg.Metadata.Get(accountIDKey))
	assert.JSONEq(t, `["READ_PUBLISHED_THREADS"]`, msg.Metadata.Get(permissionsKey))
}

func TestSessionContextMiddleware_MalformedPermissionsRejectsMessage(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handlerCalled := false

	mockHandler := func(msg *message.Message) ([]*message.Message, error) {
		handlerCalled = true
		return nil, nil
	}

	middleware := newSessionContextMiddleware(logger)
	wrappedHandler := middleware(mockHandler)

	msg := message.NewMessage("test-uuid", []byte("test payload"))
	msg.Metadata.Set(permissionsKey, "{")

	_, err := wrappedHandler(msg)

	require.Error(t, err)
	assert.False(t, handlerCalled)
}

func TestChaosDelayMiddleware_ZeroDelay(t *testing.T) {
	logger := slog.Default()
	handlerCalled := false

	mockHandler := func(msg *message.Message) ([]*message.Message, error) {
		handlerCalled = true
		return nil, nil
	}

	middleware := newChaosDelayMiddleware(0, logger)
	wrappedHandler := middleware(mockHandler)

	msg := message.NewMessage("test-uuid", []byte("test payload"))

	start := time.Now()
	_, err := wrappedHandler(msg)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.Less(t, elapsed, 10*time.Millisecond)
}

func TestChaosDelayMiddleware_WithDelay(t *testing.T) {
	logger := slog.Default()
	handlerCalled := false

	mockHandler := func(msg *message.Message) ([]*message.Message, error) {
		handlerCalled = true
		return nil, nil
	}

	maxDelay := 100 * time.Millisecond
	middleware := newChaosDelayMiddleware(maxDelay, logger)
	wrappedHandler := middleware(mockHandler)

	msg := message.NewMessage("test-uuid", []byte("test payload"))
	msg.Metadata.Set("name", "TestEvent")

	start := time.Now()
	_, err := wrappedHandler(msg)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.LessOrEqual(t, elapsed, maxDelay+10*time.Millisecond)
}

func TestChaosDelayMiddleware_MultipleCalls(t *testing.T) {
	logger := slog.Default()
	callCount := 0

	mockHandler := func(msg *message.Message) ([]*message.Message, error) {
		callCount++
		return nil, nil
	}

	maxDelay := 50 * time.Millisecond
	middleware := newChaosDelayMiddleware(maxDelay, logger)
	wrappedHandler := middleware(mockHandler)

	msg1 := message.NewMessage("test-uuid-1", []byte("test payload 1"))
	msg1.Metadata.Set("name", "TestEvent1")

	msg2 := message.NewMessage("test-uuid-2", []byte("test payload 2"))
	msg2.Metadata.Set("name", "TestEvent2")

	start := time.Now()
	_, err1 := wrappedHandler(msg1)
	_, err2 := wrappedHandler(msg2)
	elapsed := time.Since(start)

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.Equal(t, 2, callCount)

	assert.LessOrEqual(t, elapsed, 2*maxDelay+20*time.Millisecond)
}
