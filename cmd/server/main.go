package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	authImpl "trainingFinder/internal/auth"
	"trainingFinder/internal/config"
	authPkg "trainingFinder/pkg/api/auth"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
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

    httpPort := cfg.Server.HTTPPort
    ctx := context.Background()
    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
    fmt.Println("hohohoh",serverHost)
    fmt.Println("popopo",grpcPort)
    err = authPkg.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, serverHost+grpcPort, opts)
    if err != nil {
        log.Fatal("Registration error", err)
    }

    log.Printf("HTTP port %s", httpPort)
    if err := http.ListenAndServe(httpPort, mux); err != nil {
        log.Fatal("HTTP error:", err)
    }
}
