package threadapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

func Fetch(ctx context.Context, client *openapi.ClientWithResponses, mark string) (*openapi.Thread, error) {
	response, err := client.ThreadGetWithResponse(ctx, mark, &openapi.ThreadGetParams{})
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, GetError(response)
	}

	return response.JSON200, nil
}

func GetError(response *openapi.ThreadGetResponse) error {
	if response.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("thread not found")
	}

	return output.RequestErrorWithMessages("thread get request", response, response.Body, output.UnauthorizedMessage("thread get request"))
}
