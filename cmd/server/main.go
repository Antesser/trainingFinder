package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"trainingFinder/internal/app/controller"
	"trainingFinder/internal/config"

	authGRPS "trainingFinder/internal/app/auth"
	authRepository "trainingFinder/internal/repository/auth"
	authService "trainingFinder/internal/service/auth"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	dbURL := cfg.PostgresURL()
	log.Println("dbURL", dbURL)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal("Unable to create connection pool:", err)
	}
	defer pool.Close()

	authRepo := authRepository.New(pool)
	//trainingRepo := trainingRepository.New(pool)
	authSrv := authService.New(authRepo)
	//trainingSrv := trainingService.New(trainingRepo)

	// TODO:
	// Создать структуру controller в app/controller, которая будет реализовывать RunGRPC, RunHTTP и прокинуть в конструктор все grpc сервера
	// Для grpc сервера users повторить реализацию методов Register на примере auth
	// Закрепить чистую архитектуру
	grpcController := authGRPS.NewServer(authSrv)

	ctrl := controller.New(cfg.Server, grpcController)
	ctrl.Run(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	log.Println("Application ex")

}
