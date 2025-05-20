package grpc

import (
	"context"
	"github.com/alenalato/users-service/pkg/grpc"
)

func (s *UsersServer) DeleteUser(ctx context.Context, req *grpc.DeleteUserRequest) (*grpc.DeleteUserResponse, error) {
	deleteErr := s.userManager.DeleteUser(ctx, req.GetUserId())
	if deleteErr != nil {
		return nil, commonErrorToGRPCError(deleteErr)
	}

	return &grpc.DeleteUserResponse{}, nil
}
