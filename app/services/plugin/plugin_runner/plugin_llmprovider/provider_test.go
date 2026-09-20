package plugin_llmprovider

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/asset"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
)

func TestPluginProvidersReportImageInputUnsupported(t *testing.T) {
	provider := New(model_ref.NewProvider("plugin:test"), nil)

	capabilities, err := provider.ModelCapabilities(context.Background(), model_ref.ModelRef{})
	require.NoError(t, err)
	assert.Equal(t, llm_provider.CapabilitySupportUnsupported, capabilities.ImageInput)
}

func TestPluginModelRejectsImagesBeforeSendingRPC(t *testing.T) {
	provider := model_ref.NewProvider("plugin:test")
	m := &pluginModel{provider: provider, model: model_ref.NewModel("text-model")}
	req := &model.LLMRequest{Contents: []*genai.Content{{
		Role:  genai.RoleUser,
		Parts: []*genai.Part{robotresource.NewImageAssetPart(asset.AssetID(xid.New()))},
	}}}

	var gotErr error
	for _, err := range m.GenerateContent(context.Background(), req, false) {
		gotErr = err
	}
	require.Error(t, gotErr)
	assert.ErrorContains(t, gotErr, "does not support image input")
}
