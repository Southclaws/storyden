package otp

import (
	"crypto/rand"
	"fmt"
	"math"

	"github.com/Southclaws/fault"
)


func Generate() (string, error) {
	sum := make([]byte, 6)
	_, err := rand.Read(sum)
	if err != nil {
		return "", fault.Wrap(err)
	}

	value := int64(((int(sum[0]) & 0x7f) << 24) |
		((int(sum[1] & 0xff)) << 16) |
		((int(sum[2] & 0xff)) << 8) |
		(int(sum[3]) & 0xff))

	mod := int32(value % int64(math.Pow10(6)))

	return fmt.Sprintf("%06d", mod), nil
}
