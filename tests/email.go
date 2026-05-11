package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
)

func WaitForNextEmail(t *testing.T, inbox *mailer.Mock, previousCount int) mailer.MockEmail {
	t.Helper()

	require.Eventually(t, func() bool {
		return inbox.Count() > previousCount
	}, 5*time.Second, 50*time.Millisecond)

	return inbox.GetLast()
}
