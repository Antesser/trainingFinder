package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	authImpl "trainingFinder/internal/app/auth"
	"trainingFinder/internal/app/controller"
	"trainingFinder/internal/config"
	authPkg "trainingFinder/pkg/api/auth"

	authRepository "trainingFinder/internal/repository/auth"
	authService "trainingFinder/internal/service/auth"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	dbURL := cfg.PostgresURL()
	log.Println("dbURL", dbURL)
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Unable to create connection pool:", err)
	}
	defer pool.Close()

	authRepo := authRepository.New(pool)
// usersRepo := users.New(pool)   
	authSrv := authService.New(authRepo)
	//  usersSrv := users.NewService(usersRepo)


	// TODO:
	// Создать структуру controller в app/controller, которая будет реализовывать RunGRPC, RunHTTP и прокинуть в конструктор все grpc сервера
	// Для grpc сервера users повторить реализацию методов Register на примере auth
	// Закрепить чистую архитектуру


    grpcServer := grpc.NewServer()
    httpMux := runtime.NewServeMux()

    ctrl := controller.New(cfg, grpcServer, httpMux)

    if err := ctrl.RegisterServices(authSrv); err != nil {
        log.Fatal("Failed to register services:", err)
    }
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    if err := ctrl.RunGRPC(ctx); err != nil {
        log.Fatal("Failed to start gRPC server:", err)
    }
    log.Printf("gRPC server listening on %s", cfg.Server.GRPCPort)

    if err := ctrl.RunHTTP(ctx); err != nil {
        log.Fatal("Failed to start HTTP server:", err)
    }
    log.Printf("HTTP server listening on %s", cfg.Server.HTTPPort)

    с := make(chan os.Signal, 1)
    signal.Notify(с, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
    <-с

    log.Println("Shutting down gracefully...")

    cancel()

    done := make(chan struct{})
    go func() {
        ctrl.Wait()
        close(done)
    }()

    select {
    case <-done:
        log.Println("All servers stopped gracefully")
    case <-time.After(10 * time.Second):
        log.Println("Shutdown timeout exceeded, forcing exit")
    }

    log.Println("Application exited")

}
