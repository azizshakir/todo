package main

import (
	"net"

	"google.golang.org/grpc"

	"github.com/azizshakir/todo/config"
	"github.com/azizshakir/todo/service"
	"github.com/azizshakir/todo/pkg/db"
	"github.com/azizshakir/todo/pkg/logger"
	pb"github.com/azizshakir/todo/genproto"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "todo")
	defer func(l logger.Logger) {
		err := logger.Cleanup(l)
		if err != nil {
			log.Fatal("failed cleanup logger", logger.Error(err))
		}
	}(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	taskService := service.NewTaskService(connDB, log)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, taskService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
