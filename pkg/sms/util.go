package sms

import (
	c "crypto/md5"
	"encoding/hex"
)

func md5(src string) string {
	ctx := c.New()
	ctx.Write([]byte(src))
	return hex.EncodeToString(ctx.Sum(nil))
}
