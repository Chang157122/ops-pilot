package util

import (
	"crypto/md5"
	"encoding/hex"
)

// SToMD5 字符串转MD5
func SToMD5(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
