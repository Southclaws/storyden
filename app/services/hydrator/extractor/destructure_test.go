package extractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const md1 = `Check out my new oven!

![http://localhost:3000/api/v1/assets/5902d45bf0cd23d88c70b5e38652c44e2d815b08](http://localhost:3000/api/v1/assets/5902d45bf0cd23d88c70b5e38652c44e2d815b08)

Isn't it cool?

Here's a link: https://ao.com/cooking/ovens`

const md2 = `Embeds, separate line:

https://ao.com/cooking/ovens

Same line: https://x.com/southclaws bare link.

Same line, [link text](https://cla.ws) inline.
`

func TestDestructure(t *testing.T) {
	tests := []struct {
		name string
		text string
		want EnrichedProperties
	}{
		{
			name: "ovens_lol",
			text: md1,
			want: EnrichedProperties{
				Short: "Check out my new oven! http://localhost:3000/api/v1/assets/5902d45bf0cd23d88c70b5e38652c44e2d815b08 Isn't it cool? Here's a link...",
				Links: []string{
					"https://ao.com/cooking/ovens",
				},
			},
		},
		{
			name: "embedsssss",
			text: md2,
			want: EnrichedProperties{
				Short: "Embeds, separate line: Same line:   bare link. Same line,   inline.",
				Links: []string{
					"https://ao.com/cooking/ovens",
					"https://x.com/southclaws",
					"https://cla.ws",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := Destructure(tt.text)

			assert.Equal(t, tt.want.Short, ep.Short)
			assert.Equal(t, tt.want.Links, ep.Links)
		})
	}
}
