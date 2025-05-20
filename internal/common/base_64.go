package common

import (
	"encoding/base64"
	"github.com/alenalato/users-service/internal/logger"
)

func Base64Decode(str string) string {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		logger.Log.Debugf("failed to decode base64 string: %s", err.Error())

		return ""
	}

	return string(data)
}

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
