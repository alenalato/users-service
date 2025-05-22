package mongodb

import (
	"encoding/json"
	"fmt"
	"github.com/alenalato/users-service/internal/common"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/alenalato/users-service/internal/storage"
	"strconv"
	"strings"
)

const tokenDelimiter = "#"

// generateNextPageToken generates the next page token for pagination
// based on the provided filter and skip size.
// The token is a base64 encoded string that contains the serialized filter
// and the skip size, separated by a delimiter.
func generateNextPageToken(filter storage.UserFilter, skipSize int64) (string, error) {
	serializedFilter, err := serializeFilter(filter)
	if err != nil {
		return "", err
	}
	token := fmt.Sprintf("%s%s%d", serializedFilter, tokenDelimiter, skipSize)

	return common.Base64Encode(token), nil
}

// parsePageToken parses the page token and returns the user filter
// and skip size. The token is expected to be a base64 encoded string
// that contains the serialized filter and the skip size, separated by a delimiter.
func parsePageToken(
	pageToken string,
) (userFilter storage.UserFilter, skipSize int64, err error) {
	// Decode base64
	plainToken := common.Base64Decode(pageToken)
	if plainToken == "" {
		err = fmt.Errorf("cannot decode page token: %s", pageToken)
		logger.Log.Error(err)

		return storage.UserFilter{}, 0, common.NewError(err, common.ErrTypeInvalidArgument)
	}

	// Split the token
	parts := strings.Split(plainToken, tokenDelimiter)
	if len(parts) != 2 {
		err = fmt.Errorf("invalid page token format: %s", plainToken)
		logger.Log.Error(err)
	}

	filter, err := deserializeFilter(parts[0])
	if err != nil {
		return storage.UserFilter{}, 0, err
	}

	// Parse the skip size
	skipSize, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		err = fmt.Errorf("cannot parse skip size from page token: %s", err.Error())
		logger.Log.Error(err)

		return storage.UserFilter{}, 0, common.NewError(err, common.ErrTypeInvalidArgument)
	}

	return filter, skipSize, nil
}

// serializeFilter serializes the user filter into a JSON string.
func serializeFilter(filter storage.UserFilter) (string, error) {
	data, err := json.Marshal(filter)
	if err != nil {
		logger.Log.Errorf("failed to serialize filter: %v", err)

		return "", common.NewError(err, common.ErrTypeInternal)
	}
	return string(data), nil
}

// deserializeFilter deserializes the JSON string into a user filter.
func deserializeFilter(data string) (storage.UserFilter, error) {
	var filter storage.UserFilter
	err := json.Unmarshal([]byte(data), &filter)
	if err != nil {
		logger.Log.Errorf("failed to deserialize filter: %v", err)

		return storage.UserFilter{}, common.NewError(err, common.ErrTypeInvalidArgument)
	}
	return filter, nil
}
