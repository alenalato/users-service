package grpc

import (
	"context"
	"github.com/alenalato/users-service/pkg/grpc"
)

func (s *UsersServer) CreateUser(ctx context.Context, req *grpc.CreateUserRequest) (*grpc.CreateUserResponse, error) {
	user, errCreate := s.userManager.CreateUser(ctx, s.converter.fromGrpcCreateUserRequestToModel(ctx, req))
	if errCreate != nil {
		return nil, commonErrorToGRPCError(errCreate)
	}

	return &grpc.CreateUserResponse{
		User: s.converter.fromModelUserToGrpc(ctx, *user),
	}, nil
}
