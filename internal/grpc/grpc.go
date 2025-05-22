package grpc

import (
	"github.com/alenalato/users-service/internal/businesslogic"
	"github.com/alenalato/users-service/pkg/grpc"
)

// UsersServer is the server API for Users service.
type UsersServer struct {
	grpc.UnimplementedUsersServer
	// converter is used to convert between gRPC and business logic models
	converter modelConverter
	// userManager is the business logic layer for user management
	userManager businesslogic.UserManager
}

// NewUsersServer creates a new UsersServer
func NewUsersServer(userManager businesslogic.UserManager) *UsersServer {
	return &UsersServer{
		converter:   newServerModelConverter(),
		userManager: userManager,
	}
}
