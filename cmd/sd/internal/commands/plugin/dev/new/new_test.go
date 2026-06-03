package pluginnew

import "testing"

func TestNewRejectsDirectoryWithoutSlug(t *testing.T) {
	cmd := New()
	(*cmd).SetArgs([]string{"."})

	if err := (*cmd).Execute(); err == nil {
		t.Fatal("expected error")
	}
}
