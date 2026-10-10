package reaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmoji(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		in    string
		want  string
		valid bool
	}{
		{"empty", "", "", false},
		{"single byte does not panic", "a", "", false},
		{"colon at index one", "a:", "", false},
		{"plain word", "hello", "hello", false},
		{"real emoji", "👍", "👍", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := IsValidEmoji(c.in)
			assert.Equal(t, c.valid, ok)
			if c.valid {
				assert.Equal(t, c.want, got)
			}
		})
	}
}