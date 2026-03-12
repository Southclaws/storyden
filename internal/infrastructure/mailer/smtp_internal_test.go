package mailer

import "testing"

func TestUseImplicitTLS(t *testing.T) {
	if !useImplicitTLS(465) {
		t.Fatalf("expected implicit TLS on port 465")
	}

	if useImplicitTLS(587) {
		t.Fatalf("expected STARTTLS on non-implicit TLS ports")
	}
}

func TestSMTPServerAddress(t *testing.T) {
	if got := smtpServerAddress("smtp.example.com", 587); got != "smtp.example.com:587" {
		t.Fatalf("unexpected hostname address, got %q", got)
	}

	if got := smtpServerAddress("2001:db8::1", 587); got != "[2001:db8::1]:587" {
		t.Fatalf("unexpected IPv6 address, got %q", got)
	}
}
