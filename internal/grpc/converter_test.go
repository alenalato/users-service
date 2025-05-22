package grpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"testing"
	"time"

	"github.com/alenalato/users-service/internal/businesslogic"
	protogrpc "github.com/alenalato/users-service/pkg/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServerModelConverter_FromGrpcCreateUserRequestToModel(t *testing.T) {
	converter := newServerModelConverter()

	tests := []struct {
		name string
		req  *protogrpc.CreateUserRequest
		want businesslogic.UserDetails
	}{
		{
			name: "Valid request",
			req: &protogrpc.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "jdoe",
				Email:     "john.doe@example.com",
				Password:  "password123",
				Country:   "US",
			},
			want: businesslogic.UserDetails{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "jdoe",
				Email:     "john.doe@example.com",
				Password: businesslogic.PasswordDetails{
					Text: "password123",
				},
				Country: "US",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := converter.fromGrpcCreateUserRequestToModel(context.Background(), tt.req)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestServerModelConverter_FromGrpcUpdateUserRequestToModel(t *testing.T) {
	converter := newServerModelConverter()

	tests := []struct {
		name string
		req  *protogrpc.UpdateUserRequest
		want businesslogic.UserUpdate
	}{
		{
			name: "Valid request",
			req: &protogrpc.UpdateUserRequest{
				Update: &protogrpc.UpdateUserRequest_Update{
					FirstName: "Jane",
					LastName:  "Smith",
					Nickname:  "jsmith",
					Email:     "jane@smith.com",
					Country:   "UK",
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"first_name", "last_name", "nickname", "email", "country"},
				},
			},
			want: businesslogic.UserUpdate{
				FirstName:  "Jane",
				LastName:   "Smith",
				Country:    "UK",
				Nickname:   "jsmith",
				Email:      "jane@smith.com",
				UpdateMask: []string{"first_name", "last_name", "nickname", "email", "country"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := converter.fromGrpcUpdateUserRequestToModel(context.Background(), tt.req)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestServerModelConverter_FromGrpcListUsersRequestToModel(t *testing.T) {
	converter := newServerModelConverter()

	firstName := "Alice"
	lastName := "Johnson"
	country := "CA"
	tests := []struct {
		name string
		req  *protogrpc.ListUsersRequest
		want businesslogic.UserFilter
	}{
		{
			name: "Valid request with filters",
			req: &protogrpc.ListUsersRequest{
				Filter: &protogrpc.UserFilter{
					FirstName: &protogrpc.UserFilter_FirstNameFilter{
						Value: firstName,
					},
					LastName: &protogrpc.UserFilter_LastNameFilter{
						Value: lastName,
					},
					Country: &protogrpc.UserFilter_CountryFilter{
						Value: country,
					},
				},
			},
			want: businesslogic.UserFilter{
				FirstName: &firstName,
				LastName:  &lastName,
				Country:   &country,
			},
		},
		{
			name: "Valid request with partial filters",
			req: &protogrpc.ListUsersRequest{
				Filter: &protogrpc.UserFilter{
					FirstName: &protogrpc.UserFilter_FirstNameFilter{
						Value: firstName,
					},
				},
			},
			want: businesslogic.UserFilter{
				FirstName: &firstName,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := converter.fromGrpcListUsersRequestToModel(context.Background(), tt.req)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestServerModelConverter_FromModelUserToGrpc(t *testing.T) {
	converter := newServerModelConverter()

	tests := []struct {
		name string
		user businesslogic.User
		want *protogrpc.User
	}{
		{
			name: "Valid user",
			user: businesslogic.User{
				ID:        "123",
				FirstName: "Bob",
				LastName:  "Brown",
				Nickname:  "bbrown",
				Email:     "bob.brown@example.com",
				Country:   "AU",
				CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			want: &protogrpc.User{
				Id:        "123",
				FirstName: "Bob",
				LastName:  "Brown",
				Nickname:  "bbrown",
				Email:     "bob.brown@example.com",
				Country:   "AU",
				CreatedAt: timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: timestamppb.New(time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := converter.fromModelUserToGrpc(context.Background(), tt.user)
			assert.Equal(t, tt.want, got)
		})
	}
}
