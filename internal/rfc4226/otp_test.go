package rfc4226

import (
	"fmt"
	"testing"
)

func TestOTP(t *testing.T) {
	fmt.Println(OTP("12345678901234567890", 0))
}
