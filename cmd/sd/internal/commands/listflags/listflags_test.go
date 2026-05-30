package listflags

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateRejectsBadFlags(t *testing.T) {
	r := require.New(t)

	r.ErrorContains((&Flags{Page: 0, Format: FormatAuto, Output: OutputDefault}).Validate(), "--page")
	r.ErrorContains((&Flags{Page: 1, Limit: -1, Format: FormatAuto, Output: OutputDefault}).Validate(), "--limit")
	r.ErrorContains((&Flags{Page: 1, Format: "yaml", Output: OutputDefault}).Validate(), "--format")
	r.ErrorContains((&Flags{Page: 1, Format: FormatAuto, Output: "tall"}).Validate(), "--output")
}

func TestValidateAcceptsAllSupportedFormats(t *testing.T) {
	r := require.New(t)
	for _, f := range []string{FormatAuto, FormatPlain, FormatJSON, FormatJSONL} {
		r.NoError((&Flags{Page: 1, Format: f, Output: OutputDefault}).Validate(), f)
	}
}

func TestResolveFormatAutoIsPlain(t *testing.T) {
	r := require.New(t)
	r.Equal(FormatPlain, (&Flags{Format: FormatAuto}).ResolveFormat(&bytes.Buffer{}))
	r.Equal(FormatJSON, (&Flags{Format: FormatJSON}).ResolveFormat(&bytes.Buffer{}))
	r.Equal(FormatJSONL, (&Flags{Format: FormatJSONL}).ResolveFormat(&bytes.Buffer{}))
}
