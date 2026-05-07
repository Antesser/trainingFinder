package main

import (
    "context"
    "log"
    "net"
    "net/http"

    authPkg "trainingFinder/pkg/api/auth"
    authImpl "trainingFinder/internal/auth"

    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    grpcPort := ":9090"
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

    httpPort := ":8080"
    ctx := context.Background()
    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
    err := authPkg.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, "localhost"+grpcPort, opts)
    if err != nil {
        log.Fatal("Registration error", err)
    }

    log.Printf("HTTP port %s", httpPort)
    if err := http.ListenAndServe(httpPort, mux); err != nil {
        log.Fatal("HTTP error:", err)
    }
}