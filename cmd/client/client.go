package client

import (
	"context"
	"log"
	"time"

	"github.com/Antesser/trainingFinder/internal/config"
	authPkg "github.com/Antesser/trainingFinder/pkg/api/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startClient() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	conn, err := grpc.NewClient(cfg.Server.Host+cfg.Server.GRPCPort,
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
