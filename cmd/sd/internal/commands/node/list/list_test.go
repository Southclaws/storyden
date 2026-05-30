package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateVisibilities(t *testing.T) {
	r := require.New(t)

	r.NoError(validateVisibilities(nil))
	r.NoError(validateVisibilities([]string{}))
	r.NoError(validateVisibilities([]string{"draft", "review", "published", "unlisted"}))
	r.ErrorContains(validateVisibilities([]string{"private"}), "invalid --visibility: private")
	r.ErrorContains(validateVisibilities([]string{"draft", "garbage"}), "invalid --visibility: garbage")
}
