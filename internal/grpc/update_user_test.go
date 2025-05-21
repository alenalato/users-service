package grpc

import (
	"context"
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/internal/common"
	protogrpc "github.com/alenalato/users-service/pkg/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestUsersServer_UpdateUser_Success(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	req := &protogrpc.UpdateUserRequest{
		UserId: "123",
		Update: &protogrpc.UpdateUserRequest_Update{
			FirstName: "Jane",
			LastName:  "Doe",
			Country:   "Canada",
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"first_name", "last_name", "country"},
		},
	}

	userUpdate := businesslogic.UserUpdate{
		FirstName:  "Jane",
		LastName:   "Doe",
		Country:    "Canada",
		UpdateMask: []string{"first_name", "last_name", "country"},
	}

	updatedUser := &businesslogic.User{
		ID:        "123",
		FirstName: "Jane",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Email:     "john@doe.com",
		Country:   "Canada",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	grpcUser := &protogrpc.User{
		Id:        "123",
		FirstName: "Jane",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Email:     "john@doe.com",
		Country:   "Canada",
		CreatedAt: timestamppb.New(updatedUser.CreatedAt),
		UpdatedAt: timestamppb.New(updatedUser.UpdatedAt),
	}

	ts.mockConverter.EXPECT().fromGrpcUpdateUserRequestToModel(gomock.Any(), req).Return(userUpdate)

	ts.mockUserManager.EXPECT().UpdateUser(gomock.Any(), "123", userUpdate).Return(updatedUser, nil)

	ts.mockConverter.EXPECT().fromModelUserToGrpc(gomock.Any(), *updatedUser).Return(grpcUser)

	resp, err := ts.usersServer.UpdateUser(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &protogrpc.UpdateUserResponse{
		User: grpcUser,
	}, resp)
}

func TestUsersServer_UpdateUser_ManagerError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	req := &protogrpc.UpdateUserRequest{
		UserId: "123",
		Update: &protogrpc.UpdateUserRequest_Update{
			FirstName: "Jane",
			LastName:  "Doe",
			Country:   "Canada",
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"first_name", "last_name", "country"},
		},
	}

	userUpdate := businesslogic.UserUpdate{
		FirstName:  "Jane",
		LastName:   "Doe",
		Country:    "Canada",
		UpdateMask: []string{"first_name", "last_name", "country"},
	}

	ts.mockConverter.EXPECT().fromGrpcUpdateUserRequestToModel(gomock.Any(), req).Return(userUpdate)

	ts.mockUserManager.EXPECT().UpdateUser(gomock.Any(), "123", userUpdate).
		Return(nil, common.NewError(nil, common.ErrTypeNotFound))

	resp, err := ts.usersServer.UpdateUser(context.Background(), req)
	assert.Nil(t, resp)
	assert.Error(t, err)
	errGrpc, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, errGrpc.Code())
}
