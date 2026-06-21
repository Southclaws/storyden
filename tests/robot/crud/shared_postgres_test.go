package crud_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Southclaws/storyden/tests"
)

func TestMain(m *testing.M) {
	if tests.IsSharedPostgresDatabase() {
		fmt.Println("skipping robot CRUD tests on shared postgres database")
		os.Exit(0)
	}

	os.Exit(m.Run())
}
