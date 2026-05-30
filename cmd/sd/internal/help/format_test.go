package help

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTerminalReturnsFalseForNonFileWriters(t *testing.T) {
	a := assert.New(t)

	a.False(IsTerminal(&bytes.Buffer{}))
}
