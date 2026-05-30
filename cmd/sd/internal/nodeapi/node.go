package nodeapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

func Fetch(ctx context.Context, client *openapi.ClientWithResponses, slug string) (*openapi.NodeWithChildren, error) {
	response, err := client.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{})
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, GetError(response)
	}

	return response.JSON200, nil
}

func Update(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	props openapi.NodeMutableProps,
) (*openapi.NodeWithChildren, error) {
	response, err := client.NodeUpdateWithResponse(ctx, slug, props)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, UpdateError(response)
	}

	return response.JSON200, nil
}

func GetError(response *openapi.NodeGetResponse) error {
	if response.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("node not found")
	}

	return output.RequestErrorWithMessages("node get request", response, response.Body, output.UnauthorizedMessage("node get request"))
}

func UpdateError(response *openapi.NodeUpdateResponse) error {
	if response.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("node not found")
	}

	return output.RequestErrorWithMessages("node update request", response, response.Body, output.UnauthorizedMessage("node update request"))
}
