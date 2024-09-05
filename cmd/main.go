package main

import (
	"auth_service/api"
	"auth_service/api/handlers"
	"auth_service/configs"
	"auth_service/grpc"
	"auth_service/pkg/logger"
	"auth_service/storage/postgres"
	"context"
	"net"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := configs.Load()

	loggerLevel := logger.LevelDebug
	switch cfg.Environment {
	case configs.DebugMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case configs.TestMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.NewLogger(cfg.ServiceName, loggerLevel)
	defer func(l logger.ILogger) {
		logger.Cleanup(l)

	}(log)

	store, err := postgres.NewStore(context.Background(), log, cfg)
	if err != nil {
		log.Error("Failed to initialize store", logger.Error(err))
		return
	}

	serv := postgres.NewAuthRepository(store.DB, log)
	h := handlers.NewHandler(serv, log)

	r := api.Router(h)
	go func() {
		log.Info("HTTP server is running on port :8075")
		if err := r.Run(":8075"); err != nil {
			log.Error("Error while running HTTP server", logger.Error(err))
		}
	}()

	grpcServer := grpc.SetUpServer(store.Users(), log)
	lis, err := net.Listen("tcp", cfg.AuthServiceGrpcHost+cfg.AuthServiceGrpcPort)
	if err != nil {
		log.Error("Error while listening on gRPC host port", logger.Error(err))
		return
	}

	log.Info("gRPC service is running...", logger.Any("grpc port", cfg.AuthServiceGrpcPort))
	if err := grpcServer.Serve(lis); err != nil {
		log.Error("Error while running gRPC server", logger.Error(err))
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
