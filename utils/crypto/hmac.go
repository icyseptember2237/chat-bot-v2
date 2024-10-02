package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

func HMACSha1(key string, data string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}
