package grpc

import (
	"context"
	"github.com/alenalato/users-service/pkg/grpc"
)

func (s *UsersServer) ListUsers(context context.Context, req *grpc.ListUsersRequest) (*grpc.ListUsersResponse, error) {
	return nil, nil
}
