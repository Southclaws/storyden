package node_versions_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/tests"
)

func createPublishedNode(
	t *testing.T,
	root context.Context,
	cl *openapi.ClientWithResponses,
	adminSession openapi.RequestEditorFn,
	prefix string,
) *openapi.NodeCreateOK {
	t.Helper()

	published := openapi.VisibilityPublished
	name := prefix + "-" + uuid.NewString()

	node, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
		Name:       name,
		Visibility: &published,
	}, adminSession)
	tests.Ok(t, err, node)
	require.NotNil(t, node.JSON200)

	return node.JSON200
}

func createDraftVersion(
	t *testing.T,
	root context.Context,
	cl *openapi.ClientWithResponses,
	authorSession openapi.RequestEditorFn,
	nodeKey string,
	name string,
) *openapi.NodeVersion {
	t.Helper()

	create, err := cl.NodeVersionCreateWithResponse(root, nodeKey, openapi.NodeVersionCreateJSONRequestBody{
		Name: &name,
	}, authorSession)
	tests.Ok(t, err, create)
	require.NotNil(t, create.JSON200)

	return create.JSON200
}
