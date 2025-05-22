package grpc

import (
	"context"
	"github.com/alenalato/users-service/pkg/grpc"
)

// ListUsers handles the ListUsers request
func (s *UsersServer) ListUsers(ctx context.Context, req *grpc.ListUsersRequest) (*grpc.ListUsersResponse, error) {
	// Use business logic layer to list users
	users, nextPageToken, errList := s.userManager.ListUsers(
		ctx,
		// Convert gRPC request to business logic filter model
		s.converter.fromGrpcListUsersRequestToModel(ctx, req),
		int(req.GetPageSize()),
		req.GetPageToken(),
	)
	if errList != nil {
		return nil, commonErrorToGRPCError(errList)
	}

	var grpcUsers []*grpc.User
	for _, user := range users {
		// Convert business logic user back to gRPC response's User
		grpcUsers = append(grpcUsers, s.converter.fromModelUserToGrpc(ctx, user))
	}

	return &grpc.ListUsersResponse{
		Users:         grpcUsers,
		NextPageToken: nextPageToken,
	}, nil
}
