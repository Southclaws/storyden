package chat_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Southclaws/storyden/tests"
)

func TestMain(m *testing.M) {
	if tests.IsSharedPostgresDatabase() {
		fmt.Println("skipping robot chat tests on shared postgres database")
		os.Exit(0)
	}

	os.Exit(m.Run())
}
