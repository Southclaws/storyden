package tools

import "context"

func contentForAudience(ctx context.Context, content, nextAction string) (body, next *string) {
	if ToolAudienceFromContext(ctx) == ToolAudienceMCP {
		return &content, nil
	}
	return nil, &nextAction
}
