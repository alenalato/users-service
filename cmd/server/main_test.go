package main

import (
	"context"
	"fmt"
	"github.com/alenalato/users-service/internal/businesslogic/password"
	"github.com/alenalato/users-service/internal/businesslogic/user"
	"github.com/alenalato/users-service/internal/events"
	servicegrpc "github.com/alenalato/users-service/internal/grpc"
	"github.com/alenalato/users-service/internal/logger"
	"github.com/alenalato/users-service/internal/storage/mongodb"
	protogrpc "github.com/alenalato/users-service/pkg/grpc"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log"
	"net"
	"testing"
	"time"
)

var testGrpcClient protogrpc.UsersClient

const testDbName = "users-server-test"

func TestMain(m *testing.M) {
	// Initialize listener
	buffer := 101024 * 1024
	listener := bufconn.Listen(buffer)

	// Initialize MongoDB storage
	mongoDbStorage, mongoDbCloser := getTestStorage()
	// Defer closing MongoDB storage
	defer mongoDbCloser()

	// Initialize password manager
	passwordManager := password.NewBcrypt()

	// Initialize mocked Kafka event emitter
	mockEventEmitter := events.NewMockEventEmitter(gomock.NewController(nil))
	mockEventEmitter.EXPECT().EmitUserEvent(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Initialize user manager, the business logic layer
	userManager := user.NewLogic(passwordManager, mongoDbStorage, mockEventEmitter)

	// Initialize gRPC users server
	usersServer := servicegrpc.NewUsersServer(userManager)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()

	// Register users server
	protogrpc.RegisterUsersServer(grpcServer, usersServer)

	closeServer := make(chan struct{})

	// Start gRPC server asynchronously
	go func() {
		if srvErr := grpcServer.Serve(listener); srvErr != nil {
			logger.Log.Errorf("could not listen server: %v", srvErr)
			close(closeServer)
		}
	}()

	conn, err := grpc.NewClient("passthrough://",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	defer func() {
		err := listener.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}

		grpcServer.Stop()
	}()

	testGrpcClient = protogrpc.NewUsersClient(conn)

	// run tests
	m.Run()
}

func getTestStorage() (testMongoStorage *mongodb.MongoDB, closer func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct test pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pull mongodb docker image
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "4.4",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	var dbClient *mongo.Client

	// exponential backoff-retry
	err = pool.Retry(func() error {
		dbClient, err = mongodb.NewMongoDBClient(fmt.Sprintf(
			"mongodb://%s:%s",
			resource.Container.NetworkSettings.Gateway,
			resource.GetPort("27017/tcp"),
		))
		if err != nil {
			return err
		}
		return dbClient.Ping(context.TODO(), nil)
	})

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	testMongoStorage, err = mongodb.NewMongoDB(dbClient, testDbName)

	return testMongoStorage, func() {
		// kill and remove the container
		if err = pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
		// disconnect mongodb client
		if err = testMongoStorage.Close(context.TODO()); err != nil {
			panic(err)
		}
	}
}

// Test_Integration tests the integration of all server components using gRPC operations
// against an actual gRPC server instance and a real MongoDB database, event emitter is mocked for now.
// In this first version database is shared thus test cases are not isolated and can affect each other.
// It is important to structure tests in an end-to-end manner from creating users to listing them
// and finally deleting them.
func Test_Integration(t *testing.T) {
	now := time.Now().UTC()

	testUsers := make([]*protogrpc.User, 0)
	getTestUser := func(i int) *protogrpc.User {
		if i < 0 || i >= len(testUsers) {
			return nil
		}
		return testUsers[i]
	}

	createUserTests := []struct {
		name               string
		req                *protogrpc.CreateUserRequest
		expectedRes        *protogrpc.CreateUserResponse
		expectedStatusCode codes.Code
	}{
		{
			name: "Validation error - empty nickname",
			req: &protogrpc.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "",
				Password:  "password",
				Email:     "john@doe.com",
				Country:   "UK",
			},
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name: "Validation error - invalid email",
			req: &protogrpc.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe2",
				Password:  "password",
				Email:     "not-an-email",
				Country:   "UK",
			},
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name: "Validation error - empty password",
			req: &protogrpc.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe3",
				Password:  "",
				Email:     "john3@doe.com",
				Country:   "UK",
			},
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name: "Success - valid user",
			req: &protogrpc.CreateUserRequest{
				FirstName: "Alice",
				LastName:  "Smith",
				Nickname:  "alicesmith",
				Password:  "securePassword1",
				Email:     "alice@smith.com",
				Country:   "US",
			},
			expectedRes: &protogrpc.CreateUserResponse{
				User: &protogrpc.User{
					FirstName: "Alice",
					LastName:  "Smith",
					Nickname:  "alicesmith",
					Email:     "alice@smith.com",
					Country:   "US",
				},
			},
		},
		{
			name: "Already exists - duplicate nickname",
			req: &protogrpc.CreateUserRequest{
				FirstName: "Alice",
				LastName:  "Smith",
				Nickname:  "alicesmith", // same as previous
				Password:  "securePassword1",
				Email:     "alice2@smith.com",
				Country:   "US",
			},
			expectedStatusCode: codes.AlreadyExists,
		},
		{
			name: "Already exists - duplicate email",
			req: &protogrpc.CreateUserRequest{
				FirstName: "Alice",
				LastName:  "Smith",
				Nickname:  "alicesmith2",
				Password:  "securePassword1",
				Email:     "alice@smith.com", // same as previous
				Country:   "US",
			},
			expectedStatusCode: codes.AlreadyExists,
		},
		{
			name: "Success - another valid user",
			req: &protogrpc.CreateUserRequest{
				FirstName: "Bob",
				LastName:  "Johnson",
				Nickname:  "bobjohnson",
				Password:  "securePassword2",
				Email:     "bob@johnson.com",
				Country:   "CA",
			},
			expectedRes: &protogrpc.CreateUserResponse{
				User: &protogrpc.User{
					FirstName: "Bob",
					LastName:  "Johnson",
					Nickname:  "bobjohnson",
					Email:     "bob@johnson.com",
					Country:   "CA",
				},
			},
		},
		{
			name: "Success - yet another valid user",
			req: &protogrpc.CreateUserRequest{
				FirstName: "Sam",
				LastName:  "Brown",
				Nickname:  "sambrown",
				Password:  "securePassword3",
				Email:     "sam@brown.com",
				Country:   "AU",
			},
			expectedRes: &protogrpc.CreateUserResponse{
				User: &protogrpc.User{
					FirstName: "Sam",
					LastName:  "Brown",
					Nickname:  "sambrown",
					Email:     "sam@brown.com",
					Country:   "AU",
				},
			},
		},
	}
	for _, test := range createUserTests {
		t.Run(test.name, func(t *testing.T) {
			res, err := testGrpcClient.CreateUser(context.Background(), test.req)
			if test.expectedStatusCode != 0 {
				require.Error(t, err)
				errStatus, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, test.expectedStatusCode, errStatus.Code())
			} else {
				require.NoError(t, err)
				testUsers = append(testUsers, res.GetUser())
				assert.Equal(t, test.expectedRes.GetUser().GetFirstName(), res.GetUser().GetFirstName())
				assert.Equal(t, test.expectedRes.GetUser().GetLastName(), res.GetUser().GetLastName())
				assert.Equal(t, test.expectedRes.GetUser().GetNickname(), res.GetUser().GetNickname())
				assert.Equal(t, test.expectedRes.GetUser().GetEmail(), res.GetUser().GetEmail())
				assert.Equal(t, test.expectedRes.GetUser().GetCountry(), res.GetUser().GetCountry())
				assert.NotEmpty(t, res.GetUser().GetId())
				assert.WithinDuration(t, now, res.GetUser().GetCreatedAt().AsTime(), time.Minute)
				assert.WithinDuration(t, now, res.GetUser().GetUpdatedAt().AsTime(), time.Minute)
				assert.Equal(t, res.GetUser().GetCreatedAt().AsTime(), res.GetUser().GetUpdatedAt().AsTime())
			}
		})
	}

	updateUserTests := []struct {
		name               string
		req                *protogrpc.UpdateUserRequest
		expectedRes        *protogrpc.UpdateUserResponse
		expectedStatusCode codes.Code
	}{
		{
			name: "Validation error - empty update mask",
			req: &protogrpc.UpdateUserRequest{
				UserId: "someid",
				Update: &protogrpc.UpdateUserRequest_Update{
					FirstName: "Bob",
				},
				UpdateMask: nil,
			},
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name: "Not found error - user does not exist",
			req: &protogrpc.UpdateUserRequest{
				UserId: "nonexistentid",
				Update: &protogrpc.UpdateUserRequest_Update{
					FirstName: "Ghost",
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"first_name"}},
			},
			expectedStatusCode: codes.NotFound,
		},
		{
			name: "Success - update first name, last name, nickname, email, country",
			req: &protogrpc.UpdateUserRequest{
				UserId: getTestUser(0).GetId(),
				Update: &protogrpc.UpdateUserRequest_Update{
					FirstName: "Alicia",
					LastName:  "Smythe",
					Nickname:  "aliciasmythe",
					Email:     "alicia@smythe.com",
					Country:   "UK",
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"first_name", "last_name", "nickname", "email", "country"}},
			},
			expectedRes: &protogrpc.UpdateUserResponse{
				User: &protogrpc.User{
					FirstName: "Alicia",
					LastName:  "Smythe",
					Nickname:  "aliciasmythe",
					Email:     "alicia@smythe.com",
					Country:   "UK",
					CreatedAt: getTestUser(0).GetCreatedAt(),
				},
			},
		},
		{
			name: "Already exists error - duplicate nickname",
			req: &protogrpc.UpdateUserRequest{
				UserId: getTestUser(0).GetId(),
				Update: &protogrpc.UpdateUserRequest_Update{
					Nickname: "sambrown", // already exists
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"nickname"}},
			},
			expectedStatusCode: codes.AlreadyExists,
		},
		{
			name: "Already exists error - duplicate email",
			req: &protogrpc.UpdateUserRequest{
				UserId: getTestUser(0).GetId(),
				Update: &protogrpc.UpdateUserRequest_Update{
					Email: "sam@brown.com", // already exists
				},
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"email"}},
			},
			expectedStatusCode: codes.AlreadyExists,
		},
	}

	for _, test := range updateUserTests {
		t.Run(test.name, func(t *testing.T) {
			res, err := testGrpcClient.UpdateUser(context.Background(), test.req)
			if test.expectedStatusCode != 0 {
				require.Error(t, err)
				errStatus, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, test.expectedStatusCode, errStatus.Code())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedRes.GetUser().GetFirstName(), res.GetUser().GetFirstName())
				assert.Equal(t, test.expectedRes.GetUser().GetLastName(), res.GetUser().GetLastName())
				assert.Equal(t, test.expectedRes.GetUser().GetNickname(), res.GetUser().GetNickname())
				assert.Equal(t, test.expectedRes.GetUser().GetEmail(), res.GetUser().GetEmail())
				assert.Equal(t, test.expectedRes.GetUser().GetCountry(), res.GetUser().GetCountry())
				assert.Equal(t, test.expectedRes.GetUser().GetCreatedAt(), res.GetUser().GetCreatedAt())
				assert.NotEmpty(t, res.GetUser().GetId())
				assert.WithinDuration(t, now, res.GetUser().GetUpdatedAt().AsTime(), time.Minute)
			}
		})
	}

	listUsersTests := []struct {
		name               string
		req                *protogrpc.ListUsersRequest
		expectedUsersCount int
		expectedStatusCode codes.Code
		paginate           bool
	}{
		{
			name: "List all testUsers, no filter",
			req: &protogrpc.ListUsersRequest{
				PageSize: 100,
			},
			expectedUsersCount: 3,
		},
		{
			name: "List testUsers by country filter",
			req: &protogrpc.ListUsersRequest{
				Filter: &protogrpc.UserFilter{
					Country: &protogrpc.UserFilter_CountryFilter{
						Value: "CA",
					},
				},
			},
			expectedUsersCount: 1,
		},
		{
			name: "List testUsers by first name filter",
			req: &protogrpc.ListUsersRequest{
				Filter: &protogrpc.UserFilter{
					FirstName: &protogrpc.UserFilter_FirstNameFilter{
						Value: "Sam",
					},
				},
			},
			expectedUsersCount: 1,
		},
		{
			name: "List testUsers by last name filter",
			req: &protogrpc.ListUsersRequest{
				Filter: &protogrpc.UserFilter{
					LastName: &protogrpc.UserFilter_LastNameFilter{
						Value: "Johnson",
					},
				},
			},
			expectedUsersCount: 1,
		},
		{
			name: "List testUsers by combined filter",
			req: &protogrpc.ListUsersRequest{
				Filter: &protogrpc.UserFilter{
					FirstName: &protogrpc.UserFilter_FirstNameFilter{
						Value: "Alan",
					},
					LastName: &protogrpc.UserFilter_LastNameFilter{
						Value: "Johnson",
					},
				},
			},
			expectedUsersCount: 0,
		},
		{
			name: "Invalid page size (too large)",
			req: &protogrpc.ListUsersRequest{
				PageSize: 1000,
			},
			expectedStatusCode: codes.InvalidArgument,
		},
		{
			name: "Invalid page token",
			req: &protogrpc.ListUsersRequest{
				PageToken: "invalidtoken",
			},
			expectedStatusCode: codes.InvalidArgument,
		},
	}

	for _, test := range listUsersTests {
		t.Run(test.name, func(t *testing.T) {
			res, err := testGrpcClient.ListUsers(context.Background(), test.req)
			if test.expectedStatusCode != 0 {
				require.Error(t, err)
				errStatus, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, test.expectedStatusCode, errStatus.Code())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedUsersCount, len(res.GetUsers()))
			}
		})
	}

	// Test list users pagination
	t.Run("List testUsers with pagination", func(t *testing.T) {
		res, err := testGrpcClient.ListUsers(context.Background(), &protogrpc.ListUsersRequest{
			PageSize: 2,
		})
		require.NoError(t, err)
		assert.Equal(t, 2, len(res.GetUsers()))
		assert.NotEmpty(t, res.GetNextPageToken())

		res2, err := testGrpcClient.ListUsers(context.Background(), &protogrpc.ListUsersRequest{
			PageSize:  2,
			PageToken: res.GetNextPageToken(),
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(res2.GetUsers()))
		assert.Empty(t, res2.GetNextPageToken())
	})

	deleteUserTests := []struct {
		name               string
		userId             string
		expectedStatusCode codes.Code
	}{
		{
			name:   "Delete existing user",
			userId: testUsers[1].GetId(), // Bob Johnson
		},
		{
			name:               "Delete already deleted user",
			userId:             testUsers[1].GetId(),
			expectedStatusCode: codes.NotFound,
		},
		{
			name:               "Delete non-existent user",
			userId:             "nonexistentid",
			expectedStatusCode: codes.NotFound,
		},
	}

	for _, test := range deleteUserTests {
		t.Run(test.name, func(t *testing.T) {
			req := &protogrpc.DeleteUserRequest{UserId: test.userId}
			_, err := testGrpcClient.DeleteUser(context.Background(), req)
			if test.expectedStatusCode != 0 {
				require.Error(t, err)
				errStatus, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, test.expectedStatusCode, errStatus.Code())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
