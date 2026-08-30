package bindings

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func TestThemeETagMatching(t *testing.T) {
	t.Parallel()

	exact := openapi.IfNoneMatch(`"revision-a"`)
	assert.True(t, etagMatches(&exact, `"revision-a"`))

	list := openapi.IfNoneMatch(`"older", W/"revision-a"`)
	assert.True(t, etagMatches(&list, `"revision-a"`))

	other := openapi.IfNoneMatch(`"revision-b"`)
	assert.False(t, etagMatches(&other, `"revision-a"`))
	assert.False(t, etagMatches(nil, `"revision-a"`))
}
