package server

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	goose "github.com/pressly/goose/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"trainingFinder/internal/config"
	authImpl "trainingFinder/internal/controller/auth"
	authPkg "trainingFinder/pkg/api/auth"
)

func runMigrations(dbURL string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.Up(db, "./migrations"); err != nil {
		return err
	}
	log.Println("Migration done")
	return nil
}

func startServer() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	dbURL:=cfg.PostgresURL()
	log.Println("dbURL", dbURL)
	if err := runMigrations(dbURL); err != nil {
		log.Fatal("Migration failed:", err)
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Unable to create connection pool:", err)
	}
	defer pool.Close()

	serverHost := cfg.Server.Host
	grpcPort := cfg.Server.GRPCPort
	grpcServer := grpc.NewServer()
	authServer := &authImpl.Server{}
	authPkg.RegisterAuthServiceServer(grpcServer, authServer)

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
	err = authPkg.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, serverHost+grpcPort, opts)
	if err != nil {
		log.Fatal("Registration error", err)
	}
	log.Printf("http server listening on %s", cfg.Server.HTTPPort)
	if err := http.ListenAndServe(cfg.Server.HTTPPort, mux); err != nil {
		log.Fatal("HTTP error:", err)
	}

}
