package tools

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToolAudienceContext(t *testing.T) {
	ctx := context.Background()

	assert.Equal(t, ToolAudienceUnknown, ToolAudienceFromContext(ctx))
	assert.Equal(t,
		ToolAudienceRobot,
		ToolAudienceFromContext(ContextWithToolAudience(ctx, ToolAudienceRobot)),
	)
	assert.Equal(t,
		ToolAudienceMCP,
		ToolAudienceFromContext(ContextWithToolAudience(ctx, ToolAudienceMCP)),
	)
}
