package grpc

import (
	"context"
	"github.com/alenalato/users-service/pkg/grpc"
)

func (s *UsersServer) ListUsers(ctx context.Context, req *grpc.ListUsersRequest) (*grpc.ListUsersResponse, error) {
	users, nextPageToken, errList := s.userManager.ListUsers(
		ctx,
		s.converter.fromGrpcListUsersRequestToModel(ctx, req),
		int(req.GetPageSize()),
		req.GetPageToken(),
	)
	if errList != nil {
		return nil, commonErrorToGRPCError(errList)
	}

	var grpcUsers []*grpc.User
	for _, user := range users {
		grpcUsers = append(grpcUsers, s.converter.fromModelUserToGrpc(ctx, user))
	}

	return &grpc.ListUsersResponse{
		Users:         grpcUsers,
		NextPageToken: nextPageToken,
	}, nil
}
