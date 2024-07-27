package otp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateCode(t *testing.T) {
	c1, _ := Generate()
	c2, _ := Generate()
	assert.NotEqual(t, c1, c2)
}
