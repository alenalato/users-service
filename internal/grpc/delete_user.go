package grpc

import (
	"context"
	"github.com/alenalato/users-service/pkg/grpc"
)

func (s *UsersServer) DeleteUser(ctx context.Context, req *grpc.DeleteUserRequest) (*grpc.DeleteUserResponse, error) {
	errDelete := s.userManager.DeleteUser(ctx, req.GetUserId())
	if errDelete != nil {
		return nil, commonErrorToGRPCError(errDelete)
	}

	return &grpc.DeleteUserResponse{}, nil
}
