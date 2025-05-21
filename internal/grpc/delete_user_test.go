package grpc

import (
	"context"
	"github.com/alenalato/users-service/internal/common"
	protogrpc "github.com/alenalato/users-service/pkg/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestUsersServer_DeleteUser_Success(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	req := &protogrpc.DeleteUserRequest{
		UserId: "123",
	}

	ts.mockUserManager.EXPECT().DeleteUser(gomock.Any(), "123").Return(nil)

	resp, err := ts.usersServer.DeleteUser(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.IsType(t, &protogrpc.DeleteUserResponse{}, resp)
}

func TestUsersServer_DeleteUser_ManagerError(t *testing.T) {
	ts := newTestSuite(t)
	defer ts.mockCtrl.Finish()

	req := &protogrpc.DeleteUserRequest{
		UserId: "123",
	}

	ts.mockUserManager.EXPECT().DeleteUser(gomock.Any(), "123").
		Return(common.NewError(nil, common.ErrTypeNotFound))

	resp, err := ts.usersServer.DeleteUser(context.Background(), req)
	assert.Nil(t, resp)
	assert.Error(t, err)
	errGrpc, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, errGrpc.Code())
}
