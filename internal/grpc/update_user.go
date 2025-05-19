package grpc

import (
	"context"
	"github.com/alenalato/users-service/pkg/grpc"
)

func (s *UsersServer) UpdateUser(ctx context.Context, req *grpc.UpdateUserRequest) (*grpc.UpdateUserResponse, error) {
	return nil, nil
}
