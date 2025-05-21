package grpc

import (
	"context"
	"github.com/alenalato/users-service/pkg/grpc"
)

func (s *UsersServer) UpdateUser(ctx context.Context, req *grpc.UpdateUserRequest) (*grpc.UpdateUserResponse, error) {
	user, errUpdate := s.userManager.UpdateUser(
		ctx,
		req.GetUserId(),
		s.converter.fromGrpcUpdateUserRequestToModel(ctx, req),
	)
	if errUpdate != nil {
		return nil, commonErrorToGRPCError(errUpdate)
	}

	return &grpc.UpdateUserResponse{
		User: s.converter.fromModelUserToGrpc(ctx, *user),
	}, nil
}
