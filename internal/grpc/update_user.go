package grpc

import (
	"context"
	"github.com/alenalato/users-service/pkg/grpc"
)

func (s *UsersServer) UpdateUser(ctx context.Context, req *grpc.UpdateUserRequest) (*grpc.UpdateUserResponse, error) {
	user, updateErr := s.userManager.UpdateUser(
		ctx,
		req.GetUserId(),
		s.converter.fromGrpcUpdateUserRequestToModel(ctx, req),
	)
	if updateErr != nil {
		return nil, commonErrorToGRPCError(updateErr)
	}

	return &grpc.UpdateUserResponse{
		User: s.converter.fromModelUserToGrpc(ctx, *user),
	}, nil
}
