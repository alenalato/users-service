package grpc

import (
	"context"
	"github.com/alenalato/users-service/pkg/grpc"
)

// CreateUser handles the CreateUser request
func (s *UsersServer) CreateUser(ctx context.Context, req *grpc.CreateUserRequest) (*grpc.CreateUserResponse, error) {
	// Use business logic layer to create a new user
	user, errCreate := s.userManager.CreateUser(
		ctx,
		// Convert gRPC request to business logic details model
		s.converter.fromGrpcCreateUserRequestToModel(ctx, req),
	)
	if errCreate != nil {
		return nil, commonErrorToGRPCError(errCreate)
	}

	return &grpc.CreateUserResponse{
		// Convert business logic user back to gRPC response's User
		User: s.converter.fromModelUserToGrpc(ctx, *user),
	}, nil
}
