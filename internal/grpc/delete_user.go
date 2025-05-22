package grpc

import (
	"context"
	"github.com/alenalato/users-service/pkg/grpc"
)

// DeleteUser handles the DeleteUser request
func (s *UsersServer) DeleteUser(ctx context.Context, req *grpc.DeleteUserRequest) (*grpc.DeleteUserResponse, error) {
	// Use business logic layer to delete a user
	errDelete := s.userManager.DeleteUser(ctx, req.GetUserId())
	if errDelete != nil {
		return nil, commonErrorToGRPCError(errDelete)
	}

	return &grpc.DeleteUserResponse{}, nil
}
