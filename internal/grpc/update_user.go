package grpc

import (
	"context"
	"github.com/alenalato/users-service/pkg/grpc"
)

func (s *UsersServer) UpdateUser(ctx context.Context, req *grpc.UpdateUserRequest) (*grpc.UpdateUserResponse, error) {
	// Use business logic layer to update a user
	user, errUpdate := s.userManager.UpdateUser(
		ctx,
		req.GetUserId(),
		// Convert gRPC request to business logic update model
		s.converter.fromGrpcUpdateUserRequestToModel(ctx, req),
	)
	if errUpdate != nil {
		return nil, commonErrorToGRPCError(errUpdate)
	}

	return &grpc.UpdateUserResponse{
		// Convert business logic user back to gRPC response's User
		User: s.converter.fromModelUserToGrpc(ctx, *user),
	}, nil
}
