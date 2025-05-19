package grpc

import (
	"context"
	"github.com/alenalato/users-service/pkg/grpc"
)

func (s *UsersServer) CreateUser(ctx context.Context, req *grpc.CreateUserRequest) (*grpc.CreateUserResponse, error) {
	user, createErr := s.userManager.CreateUser(ctx, s.converter.fromGrpcCreateUserRequestToModel(ctx, req))
	if createErr != nil {
		return nil, createErr
	}

	return &grpc.CreateUserResponse{
		User: s.converter.fromModelUserToGrpc(ctx, *user),
	}, nil
}
