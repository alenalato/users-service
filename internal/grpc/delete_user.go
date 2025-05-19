package grpc

import (
	"context"
	"github.com/alenalato/users-service/pkg/grpc"
)

func (s *UsersServer) DeleteUser(ctx context.Context, req *grpc.DeleteUserRequest) (*grpc.DeleteUserResponse, error) {
	return nil, nil
}
