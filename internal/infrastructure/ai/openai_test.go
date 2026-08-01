package ai

import (
	"testing"

	"github.com/openai/openai-go/v3/responses"
	"github.com/stretchr/testify/assert"
)

func TestOpenAIResponseError(t *testing.T) {
	assert.ErrorContains(t, openAIResponseError(responses.Response{
		Status: responses.ResponseStatusFailed,
		Error:  responses.ResponseError{Code: responses.ResponseErrorCodeServerError, Message: "upstream unavailable"},
	}), "server_error: upstream unavailable")

	assert.ErrorContains(t, openAIResponseError(responses.Response{
		Status:            responses.ResponseStatusIncomplete,
		IncompleteDetails: responses.ResponseIncompleteDetails{Reason: "max_output_tokens"},
	}), "incomplete: max_output_tokens")
}
