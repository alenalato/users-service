package grpc

import (
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/pkg/grpc"
)

type UsersServer struct {
	grpc.UnimplementedUsersServer
	converter   modelConverter
	userManager businesslogic.UserManager
}

func NewUsersServer(userManager businesslogic.UserManager) *UsersServer {
	return &UsersServer{
		converter:   newServerModelConverter(),
		userManager: userManager,
	}
}
