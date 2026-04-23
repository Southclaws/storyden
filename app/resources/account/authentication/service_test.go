package authentication

import "testing"

func TestNewServiceCanonicalisesBuiltIn(t *testing.T) {
	got := NewService(ServiceOAuthGitHub.String())

	if got != AuthServiceOAuthGitHub {
		t.Fatalf("expected built-in service to equal canonical alias, got=%#v want=%#v", got, AuthServiceOAuthGitHub)
	}

	if !got.IsBuiltIn() {
		t.Fatal("expected built-in service")
	}

	bis, ok := got.BuiltIn()
	if !ok {
		t.Fatal("expected built-in value")
	}

	if *bis != ServiceOAuthGitHub {
		t.Fatalf("expected built-in value %q got %q", ServiceOAuthGitHub, *bis)
	}
}

func TestNewServiceAcceptsCustomValues(t *testing.T) {
	custom := "oauth_custom_plugin"
	got := NewService(custom)

	if got.IsBuiltIn() {
		t.Fatal("expected custom service to not be built-in")
	}

	if _, ok := got.BuiltIn(); ok {
		t.Fatal("expected no built-in value for custom service")
	}

	if got.String() != custom {
		t.Fatalf("expected custom string %q got %q", custom, got.String())
	}
}
