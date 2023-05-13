package phone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateCode(t *testing.T) {
	c1, _ := generateCode()
	c2, _ := generateCode()
	assert.NotEqual(t, c1, c2)
}
