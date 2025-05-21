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
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestUsersServer_CreateUser_Success(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	req := &protogrpc.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Password:  "password123",
		Email:     "john@doe.com",
		Country:   "USA",
	}

	userDetails := businesslogic.UserDetails{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Email:     "john@doe.com",
		Password: businesslogic.PasswordDetails{
			Text: "password123",
		},
		Country: "USA",
	}

	user := &businesslogic.User{
		ID:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Email:     "john@doe.com",
		Country:   "USA",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	grpcUser := &protogrpc.User{
		Id:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Email:     "john@doe.com",
		Country:   "USA",
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}

	ts.mockConverter.EXPECT().fromGrpcCreateUserRequestToModel(gomock.Any(), req).Return(userDetails)

	ts.mockUserManager.EXPECT().CreateUser(gomock.Any(), userDetails).Return(user, nil)

	ts.mockConverter.EXPECT().fromModelUserToGrpc(gomock.Any(), *user).Return(grpcUser)

	resp, err := ts.usersServer.CreateUser(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &protogrpc.CreateUserResponse{
		User: grpcUser,
	}, resp)
}

func TestUsersServer_CreateUser_ManagerError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	req := &protogrpc.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Password:  "password123",
		Email:     "john@doe.com",
		Country:   "USA",
	}

	userDetails := businesslogic.UserDetails{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Email:     "john@doe.com",
		Password: businesslogic.PasswordDetails{
			Text: "password123",
		},
		Country: "USA",
	}

	ts.mockConverter.EXPECT().fromGrpcCreateUserRequestToModel(gomock.Any(), req).Return(userDetails)

	ts.mockUserManager.EXPECT().CreateUser(gomock.Any(), userDetails).
		Return(nil, common.NewError(nil, common.ErrTypeAlreadyExists))

	resp, err := ts.usersServer.CreateUser(context.Background(), req)
	assert.Nil(t, resp)
	assert.Error(t, err)
	errGrpc, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, errGrpc.Code())
}
