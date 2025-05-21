package mongodb

import (
	"errors"
	"github.com/alenalato/users-service/internal/common"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/alenalato/users-service/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestSerializeFilter(t *testing.T) {
	firstName := "John"
	lastName := "Doe"
	country := "USA"

	tests := []struct {
		name          string
		input         storage.UserFilter
		expectError   bool
		expectedError error
	}{
		{
			name: "Success",
			input: storage.UserFilter{
				FirstName: &firstName,
				LastName:  &lastName,
				Country:   &country,
			},
			expectError: false,
		},
		// not reproducible
		//{
		//	name: "Failure on marshalling",
		//	input: storage.UserFilter{},
		//	expectError:   true,
		//	expectedError: common.NewError(nil, common.ErrTypeInternal),
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := serializeFilter(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.True(t, errors.Is(err, tt.expectedError))
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)

				// Verify deserialization for valid cases
				if !tt.expectError {
					deserializedFilter, err := deserializeFilter(result)
					assert.NoError(t, err)
					assert.Equal(t, tt.input, deserializedFilter)
				}
			}
		})
	}
}

func TestDeserializeFilter(t *testing.T) {
	firstName := "John"
	lastName := "Doe"
	country := "USA"

	tests := []struct {
		name          string
		input         string
		expectError   bool
		expectedError common.ErrorType
	}{
		{
			name:        "Success",
			input:       `{"first_name":"John","last_name":"Doe","country":"USA"}`,
			expectError: false,
		},
		{
			name:          "Failure on unmarshalling",
			input:         `invalid_json`,
			expectError:   true,
			expectedError: common.ErrTypeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := deserializeFilter(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != 0 {
					assert.Error(t, err)
					var commonErr common.Error
					require.ErrorAs(t, err, &commonErr)
					assert.Equal(t, tt.expectedError, err.(common.Error).Type())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, firstName, *result.FirstName)
				assert.Equal(t, lastName, *result.LastName)
				assert.Equal(t, country, *result.Country)
			}
		})
	}
}

func TestGenerateNextPageToken(t *testing.T) {
	firstName := "John"
	lastName := "Doe"
	country := "USA"

	type args struct {
		filter   storage.UserFilter
		skipSize int64
	}
	tests := []struct {
		name          string
		input         args
		expectedToken string
		expectError   bool
		expectedError common.ErrorType
	}{
		{
			name: "Success",
			input: args{
				filter: storage.UserFilter{
					FirstName: &firstName,
					LastName:  &lastName,
					Country:   &country,
				},
				skipSize: 10,
			},
			expectError: false,
		},
		// not reproducible
		//{
		//	name:          "Failure on generating token",
		//	input:         args{},
		//	expectError:   true,
		//	expectedError: common.ErrTypeInternal,
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generateNextPageToken(tt.input.filter, tt.input.skipSize)
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != 0 {
					assert.Error(t, err)
					var commonErr common.Error
					require.ErrorAs(t, err, &commonErr)
					assert.Equal(t, tt.expectedError, err.(common.Error).Type())
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)

				// verify deserialization for valid cases
				deserializedFilter, skipSize, err := parsePageToken(result)
				assert.NoError(t, err)
				assert.Equal(t, tt.input.filter, deserializedFilter)
				assert.Equal(t, tt.input.skipSize, skipSize)
			}
		})
	}
}

func TestParsePageToken(t *testing.T) {
	firstName := "John"
	lastName := "Doe"
	country := "USA"

	tests := []struct {
		name          string
		input         string
		expectError   bool
		expectedError common.ErrorType
	}{
		{
			name:        "Success",
			input:       common.Base64Encode(`{"first_name":"John","last_name":"Doe","country":"USA"}#10`),
			expectError: false,
		},
		{
			name:          "Failure on invalid base64",
			input:         "invalid_base64",
			expectError:   true,
			expectedError: common.ErrTypeInternal,
		},
		{
			name:          "Failure on invalid token format",
			input:         common.Base64Encode("invalid_token_format"),
			expectError:   true,
			expectedError: common.ErrTypeInternal,
		},
		{
			name:          "Failure on invalid filter serialization",
			input:         common.Base64Encode(`invalid_filter#10`),
			expectError:   true,
			expectedError: common.ErrTypeInternal,
		},
		{
			name:          "Failure on invalid skip size",
			input:         common.Base64Encode(`{"first_name":"John","last_name":"Doe","country":"USA"}#invalid_skip_size`),
			expectError:   true,
			expectedError: common.ErrTypeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, skipSize, err := parsePageToken(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != 0 {
					assert.Error(t, err)
					var commonErr common.Error
					require.ErrorAs(t, err, &commonErr)
					assert.Equal(t, tt.expectedError, err.(common.Error).Type())
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
				assert.Equal(t, int64(10), skipSize)

				assert.Equal(t, firstName, *result.FirstName)
				assert.Equal(t, lastName, *result.LastName)
				assert.Equal(t, country, *result.Country)
			}
		})
	}
}
