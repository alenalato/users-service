package grpc

import (
	"errors"
	"testing"

	"github.com/alenalato/users-service/internal/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCommonErrorToGRPCError(t *testing.T) {
	tests := []struct {
		name         string
		inputError   error
		expectedCode codes.Code
	}{
		{
			name:         "NotFound error",
			inputError:   common.NewError(errors.New("resource not found"), common.ErrTypeNotFound),
			expectedCode: codes.NotFound,
		},
		{
			name:         "AlreadyExists error",
			inputError:   common.NewError(errors.New("resource already exists"), common.ErrTypeAlreadyExists),
			expectedCode: codes.AlreadyExists,
		},
		{
			name:         "InvalidArgument error",
			inputError:   common.NewError(errors.New("invalid argument"), common.ErrTypeInvalidArgument),
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "Internal error",
			inputError:   common.NewError(errors.New("internal error"), common.ErrTypeInternal),
			expectedCode: codes.Internal,
		},
		{
			name:         "Unknown error type",
			inputError:   common.NewError(errors.New("unknown error"), common.ErrTypeUnknown),
			expectedCode: codes.Internal,
		},
		{
			name:         "Non-common error",
			inputError:   errors.New("some other error"),
			expectedCode: codes.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := commonErrorToGRPCError(tt.inputError)
			st, _ := status.FromError(err)

			if st.Code() != tt.expectedCode {
				t.Errorf("expected code %v, got %v", tt.expectedCode, st.Code())
			}
		})
	}
}
