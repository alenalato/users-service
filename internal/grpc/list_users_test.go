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

func TestUsersServer_ListUsers_Success(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	firstName := "John"

	req := &protogrpc.ListUsersRequest{
		Filter: &protogrpc.UserFilter{
			FirstName: &protogrpc.UserFilter_FirstNameFilter{
				Value: firstName,
			},
		},
		PageSize:  10,
		PageToken: "token123",
	}

	userFilter := businesslogic.UserFilter{
		FirstName: &firstName,
	}

	users := []businesslogic.User{
		{
			ID:        "123",
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "johndoe",
			Email:     "john@doe.com",
			Country:   "USA",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "456",
			FirstName: "Jane",
			LastName:  "Smith",
			Nickname:  "janesmith",
			Email:     "jane@smith.com",
			Country:   "Canada",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	grpcUsers := []*protogrpc.User{
		{
			Id:        "123",
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "johndoe",
			Email:     "john@doe.com",
			Country:   "USA",
			CreatedAt: timestamppb.New(users[0].CreatedAt),
			UpdatedAt: timestamppb.New(users[0].UpdatedAt),
		},
		{
			Id:        "456",
			FirstName: "Jane",
			LastName:  "Smith",
			Nickname:  "janesmith",
			Email:     "john@doe.com",
			Country:   "Canada",
			CreatedAt: timestamppb.New(users[1].CreatedAt),
			UpdatedAt: timestamppb.New(users[1].UpdatedAt),
		},
	}

	ts.mockConverter.EXPECT().fromGrpcListUsersRequestToModel(gomock.Any(), req).Return(userFilter)

	ts.mockUserManager.EXPECT().ListUsers(gomock.Any(), userFilter, int(req.PageSize), req.PageToken).
		Return(users, "nextToken123", nil)

	ts.mockConverter.EXPECT().fromModelUserToGrpc(gomock.Any(), users[0]).Return(grpcUsers[0])
	ts.mockConverter.EXPECT().fromModelUserToGrpc(gomock.Any(), users[1]).Return(grpcUsers[1])

	resp, err := ts.usersServer.ListUsers(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &protogrpc.ListUsersResponse{
		Users:         grpcUsers,
		NextPageToken: "nextToken123",
	}, resp)
}

func TestUsersServer_ListUsers_ManagerError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	firstName := "John"

	req := &protogrpc.ListUsersRequest{
		Filter: &protogrpc.UserFilter{
			FirstName: &protogrpc.UserFilter_FirstNameFilter{
				Value: firstName,
			},
		},
		PageSize:  10,
		PageToken: "token123",
	}

	userFilter := businesslogic.UserFilter{
		FirstName: &firstName,
	}

	ts.mockConverter.EXPECT().fromGrpcListUsersRequestToModel(gomock.Any(), req).Return(userFilter)

	ts.mockUserManager.EXPECT().ListUsers(gomock.Any(), userFilter, int(req.PageSize), req.PageToken).
		Return(nil, "", common.NewError(nil, common.ErrTypeInternal))

	resp, err := ts.usersServer.ListUsers(context.Background(), req)
	assert.Nil(t, resp)
	assert.Error(t, err)
	errGrpc, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, errGrpc.Code())
}
