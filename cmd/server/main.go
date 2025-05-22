package main

import (
	"context"
	"fmt"
	"github.com/alenalato/users-service/internal/businesslogic/password"
	"github.com/alenalato/users-service/internal/businesslogic/user"
	"github.com/alenalato/users-service/internal/events/kafka"
	"github.com/alenalato/users-service/internal/storage/mongodb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	servicegrpc "github.com/alenalato/users-service/internal/grpc"
	"github.com/alenalato/users-service/internal/logger"
	protogrpc "github.com/alenalato/users-service/pkg/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	defer func(Log *zap.SugaredLogger) {
		_ = Log.Sync()
	}(logger.Log)

	ctx := context.Background()

	grpcListenAddress := fmt.Sprintf(
		"%s:%s",
		os.Getenv("GRPC_LISTEN_HOST"),
		os.Getenv("GRPC_LISTEN_PORT"),
	)

	// Initialize TCP listener
	listener, err := net.Listen("tcp", grpcListenAddress)
	if err != nil {
		logger.Log.Fatalf("could not listen: %v", err)
	}
	logger.Log.Infof("TCP listener initialized on %s", grpcListenAddress)

	// Initialize MongoDB storage
	mongoDbStorage, mongodbErr := mongodb.NewMongoDB(
		nil,
		os.Getenv("MONGODB_DATABASE"),
	)
	if mongodbErr != nil {
		logger.Log.Fatalf("could not initialize MongoDB storage: %v", mongodbErr)
	} else {
		logger.Log.Infof("MongoDB storage initialized")
	}
	// Defer closing MongoDB storage
	defer func(mongoDbStorage *mongodb.MongoDB, ctx context.Context) {
		err := mongoDbStorage.Close(ctx)
		if err != nil {
			logger.Log.Errorf("could not close MongoDB storage: %v", err)
		}
	}(mongoDbStorage, ctx)

	// Initialize password manager
	passwordManager := password.NewBcrypt()

	// Initialize Kafka event emitter
	kafkaEventEmitter, kafkaErr := kafka.NewEventEmitter(
		os.Getenv("KAFKA_EVENT_EMITTER_TOPIC_NAME"),
		kafka.Config{
			Addresses: strings.Split(os.Getenv("KAFKA_ADDRESSES"), ","),
		},
	)
	if kafkaErr != nil {
		logger.Log.Fatalf("could not initialize kafka emitter: %v", kafkaErr)
	} else {
		logger.Log.Infof("Kafka event emitter initialized")
	}
	// Defer closing Kafka event emitter
	defer func(kafkaEventEmitter *kafka.EventEmitter) {
		err := kafkaEventEmitter.Close()
		if err != nil {
			logger.Log.Errorf("could not close kafka emitter: %v", err)
		}
	}(kafkaEventEmitter)

	// Initialize user manager, the business logic layer
	userManager := user.NewLogic(passwordManager, mongoDbStorage, kafkaEventEmitter)

	// Initialize gRPC users server
	usersServer := servicegrpc.NewUsersServer(userManager)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()

	// Register users server
	protogrpc.RegisterUsersServer(grpcServer, usersServer)
	reflection.Register(grpcServer)

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	closeServer := make(chan struct{})

	// Start gRPC server asynchronously
	go func() {
		if srvErr := grpcServer.Serve(listener); srvErr != nil {
			logger.Log.Errorf("could not listen GRPC(%s): %v", grpcListenAddress, srvErr)
			close(closeServer)
		}
	}()

	// Handle shutdown signals
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGTERM, os.Interrupt)
		<-sigint
		close(closeServer)
	}()

	<-closeServer

	logger.Log.Infof("Waiting for gRPC server to close")

	grpcServer.GracefulStop()
}
