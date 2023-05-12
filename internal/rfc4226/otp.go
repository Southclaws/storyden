package rfc4226

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"time"
)

func OTP(shared string, count int) string {
	payload := fmt.Sprintf("%s%d", shared, time.Now().UTC().UnixMilli())

	h := hmac.New(sha1.New, []byte(payload))
	sum := h.Sum(nil)

	offset := sum[len(sum)-1] & 0xf
	truncatedHash := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	finalOTP := (truncatedHash % (10 ^ 6))

	return fmt.Sprintf("%06d", finalOTP)
}
