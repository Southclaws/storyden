package tools

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContentForAudience(t *testing.T) {
	body, next := contentForAudience(ContextWithToolAudience(context.Background(), ToolAudienceMCP), "body", "open it")
	require.NotNil(t, body)
	assert.Equal(t, "body", *body)
	assert.Nil(t, next)

	body, next = contentForAudience(ContextWithToolAudience(context.Background(), ToolAudienceRobot), "body", "open it")
	assert.Nil(t, body)
	require.NotNil(t, next)
	assert.Equal(t, "open it", *next)

	body, next = contentForAudience(context.Background(), "body", "open it")
	assert.Nil(t, body)
	require.NotNil(t, next)
	assert.Equal(t, "open it", *next)
}
