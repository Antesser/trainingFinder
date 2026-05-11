package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	authImpl "trainingFinder/internal/app/auth"
	"trainingFinder/internal/config"
	authPkg "trainingFinder/pkg/api/auth"

	authRepository "trainingFinder/internal/repository/auth"
	authService "trainingFinder/internal/service/auth"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	authSrv := authService.New(authRepo)

	serverHost := cfg.Server.Host
	grpcPort := cfg.Server.GRPCPort
	grpcServer := grpc.NewServer()

	// TODO:
	// Создать структуру controller в app/controller, которая будет реализовывать RunGRPC, RunHTTP и прокинуть в конструктор все grpc сервера
	// Для grpc сервера users повторить реализацию методов Register на примере auth
	// Закрепить чистую архитектуру

	go func() {
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatal("failure with listen:", err)
		}
		log.Printf("grpc server listening on %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("grpc server error:", err)
		}
	}()

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	log.Printf("http server listening on %s", cfg.Server.HTTPPort)
	if err := http.ListenAndServe(cfg.Server.HTTPPort, mux); err != nil {
		log.Fatal("HTTP error:", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-c

	ctx, shutdown := context.WithTimeout(ctx, 5*time.Second)
	defer shutdown()

	// shutdown
}
