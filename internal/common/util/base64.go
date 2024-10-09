package util

import (
	"encoding/base64"
	"opsPilot/internal/pkg/log"
)

func StrToBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64ToStr(base64Str string) string {
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		log.Logger.Errorf("decoding base64 failed! err: %v", err)
		return ""
	}
	return string(data)
}
