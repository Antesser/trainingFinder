package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    authPkg "trainingFinder/pkg/api/auth"
)

func main() {
    conn, err := grpc.NewClient("localhost:9090",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatal("failed to connect: ", err)
    }
    defer conn.Close()

    client := authPkg.NewAuthServiceClient(conn)

    req := &authPkg.SignUpRequest{
        Login:    "Lex",
        Password: "Pex",
    }

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    resp, err := client.SignUp(ctx, req)
    if err != nil {
        log.Fatal("SignUp issues: ", err)
    }

    log.Printf("AccessToken: %s", resp.AccessToken)
    log.Printf("RefreshToken: %s", resp.RefreshToken)
}