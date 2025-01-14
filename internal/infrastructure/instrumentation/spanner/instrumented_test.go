package spanner

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getLoc(t *testing.T) {
	a := assert.New(t)

	fn := func() {
		pkg, fn, loc := getLoc(runtime.Caller(1))
		a.Equal("spanner", pkg)
		a.Equal("spanner.Test_getLoc", fn)
		a.Contains(loc, "instrumented_test.go")
	}
	fn()
}
