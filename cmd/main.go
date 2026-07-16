package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	user "github.com/fnuritdinov/proto/userpb"
	"github.com/fnuritdinov/user-service/internal/config"
	"github.com/fnuritdinov/user-service/internal/repository"
	"github.com/fnuritdinov/user-service/internal/server"
	"github.com/fnuritdinov/user-service/internal/service"
	"github.com/fnuritdinov/user-service/pkg/cache"
	"github.com/fnuritdinov/user-service/pkg/db"
	"github.com/fnuritdinov/user-service/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	cfg, err := config.New("./config/config.env")
	if err != nil {
		log.Fatal("config.New", err)
	}

	lg, err := logger.New(true)
	if err != nil {
		log.Fatal("failed to create logger", err)
	}

	conn, err := db.New(db.Option{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		lg.Error("failed to connect to db:", zap.Error(err))
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis, err := net.Listen(cfg.NETWORK, cfg.ADDRESS)
	if err != nil {
		lg.Error("failed to listen: %v", zap.Error(err))
	}

	cch, err := cache.New(ctx, cfg.ADDRESS)
	if err != nil {
		lg.Error("failed to cache: &v", zap.Error(err))
	}
	grpcServer := grpc.NewServer()

	userRepo := repository.New(conn)

	userService := service.New(userRepo, cch)

	userServer := server.New(*lg, userService)

	user.RegisterUserServiceServer(grpcServer, userServer)

	reflection.Register(grpcServer)

	go func() {
		lg.Info("server listening at %v", zap.String("addr", lis.Addr().String()))
		if err = grpcServer.Serve(lis); err != nil {
			lg.Error("failed to serve: %v", zap.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	lg.Info("shutting down server...")
	grpcServer.GracefulStop()
	lg.Info("server stopped")
}
