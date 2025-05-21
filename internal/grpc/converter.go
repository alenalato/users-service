package grpc

import (
	"context"
	"github.com/alenalato/users-service/internal/businesslogic"
	protogrpc "github.com/alenalato/users-service/pkg/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockgen -destination=converter_mock.go -package=grpc github.com/alenalato/users-service/internal/grpc modelConverter

type modelConverter interface {
	fromGrpcCreateUserRequestToModel(ctx context.Context, req *protogrpc.CreateUserRequest) businesslogic.UserDetails
	fromGrpcUpdateUserRequestToModel(ctx context.Context, req *protogrpc.UpdateUserRequest) businesslogic.UserUpdate
	fromGrpcListUsersRequestToModel(ctx context.Context, req *protogrpc.ListUsersRequest) businesslogic.UserFilter
	fromModelUserToGrpc(ctx context.Context, user businesslogic.User) *protogrpc.User
}

type serverModelConverter struct{}

var _ modelConverter = new(serverModelConverter)

func (c *serverModelConverter) fromGrpcCreateUserRequestToModel(_ context.Context, req *protogrpc.CreateUserRequest) businesslogic.UserDetails {
	return businesslogic.UserDetails{
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Nickname:  req.GetNickname(),
		Email:     req.GetEmail(),
		Password: businesslogic.PasswordDetails{
			Text: req.GetPassword(),
		},
		Country: req.GetCountry(),
	}
}

func (c *serverModelConverter) fromGrpcUpdateUserRequestToModel(_ context.Context, req *protogrpc.UpdateUserRequest) businesslogic.UserUpdate {
	userUpdate := businesslogic.UserUpdate{
		FirstName:  req.GetUpdate().GetFirstName(),
		LastName:   req.GetUpdate().GetLastName(),
		Country:    req.GetUpdate().GetCountry(),
		UpdateMask: req.GetUpdateMask().GetPaths(),
	}

	return userUpdate
}

func (c *serverModelConverter) fromGrpcListUsersRequestToModel(_ context.Context, req *protogrpc.ListUsersRequest) businesslogic.UserFilter {
	userFilter := businesslogic.UserFilter{}

	if req.GetFilter().GetFirstName() != nil {
		firstName := req.GetFilter().GetFirstName().GetValue()
		userFilter.FirstName = &firstName
	}
	if req.GetFilter().GetLastName() != nil {
		lastName := req.GetFilter().GetLastName().GetValue()
		userFilter.LastName = &lastName
	}
	if req.GetFilter().GetCountry() != nil {
		country := req.GetFilter().GetCountry().GetValue()
		userFilter.Country = &country
	}

	return userFilter
}

func (c *serverModelConverter) fromModelUserToGrpc(_ context.Context, user businesslogic.User) *protogrpc.User {
	return &protogrpc.User{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Country:   user.Country,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func newServerModelConverter() *serverModelConverter {
	return &serverModelConverter{}
}
