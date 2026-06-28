package semdexer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
)

func TestNewSemdexerDoesNotResolveEmbedderAtBoot(t *testing.T) {
	t.Parallel()

	_, err := newSemdexer(
		context.Background(),
		config.Config{
			SemdexProvider:  "chromem",
			SemdexLocalPath: t.TempDir(),
		},
		nil,
		nil,
		nil,
	)

	require.NoError(t, err)
}
