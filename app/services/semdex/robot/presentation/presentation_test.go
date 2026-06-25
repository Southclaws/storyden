package presentation

import (
	"testing"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestParsePresentationMarkup(t *testing.T) {
	t.Parallel()

	nodeID := utils.Must(xid.FromString("cto7n8ifunp55p1bujv0"))
	profileID := utils.Must(xid.FromString("cto7nm2funp55p1bujvg"))

	tests := []struct {
		name  string
		input string
		want  []Part
	}{
		{
			name:  "plain text",
			input: "Hello world.",
			want:  []Part{{Kind: PartText, Text: "Hello world."}},
		},
		{
			name:  "standalone markdown SDR link preserves order",
			input: "Hello.\n\n[Documentation Hub](sdr:node/cto7n8ifunp55p1bujv0)\n\nWorld.",
			want: []Part{
				{Kind: PartText, Text: "Hello.\n\n"},
				{Kind: PartRenderCard, Ref: &datagraph.Ref{Kind: datagraph.KindNode, ID: nodeID}},
				{Kind: PartText, Text: "\n\nWorld."},
			},
		},
		{
			name:  "standalone html anchor SDR link preserves order",
			input: "Hello.\n\n<a href=\"sdr:profile/cto7nm2funp55p1bujvg\">@southclaws</a>\n\nWorld.",
			want: []Part{
				{Kind: PartText, Text: "Hello.\n\n"},
				{Kind: PartRenderCard, Ref: &datagraph.Ref{Kind: datagraph.KindProfile, ID: profileID}},
				{Kind: PartText, Text: "\n\nWorld."},
			},
		},
		{
			name:  "multiple standalone cards",
			input: "[Documentation Hub](sdr:node/cto7n8ifunp55p1bujv0)\n\n[@southclaws](sdr:profile/cto7nm2funp55p1bujvg)",
			want: []Part{
				{Kind: PartRenderCard, Ref: &datagraph.Ref{Kind: datagraph.KindNode, ID: nodeID}},
				{Kind: PartText, Text: "\n\n"},
				{Kind: PartRenderCard, Ref: &datagraph.Ref{Kind: datagraph.KindProfile, ID: profileID}},
			},
		},
		{
			name:  "inline SDR link remains text",
			input: "This is [Documentation Hub](sdr:node/cto7n8ifunp55p1bujv0), it is what you need.",
			want:  []Part{{Kind: PartText, Text: "This is [Documentation Hub](sdr:node/cto7n8ifunp55p1bujv0), it is what you need."}},
		},
		{
			name:  "unknown html is preserved as text",
			input: "Hello <strong>world</strong>.",
			want:  []Part{{Kind: PartText, Text: "Hello <strong>world</strong>."}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, Parse(tt.input))
		})
	}
}
