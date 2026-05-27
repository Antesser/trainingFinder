package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"trainingFinder/internal/app/controller"
	"trainingFinder/internal/config"
	your_topic_name "trainingFinder/internal/kafka/producer/your-topic-name"

	trainingGRPS "trainingFinder/internal/app/training"
	userGRPS "trainingFinder/internal/app/users"
	trainingRepository "trainingFinder/internal/repository/training"
	userRepository "trainingFinder/internal/repository/user"
	trainingService "trainingFinder/internal/service/training"
	userService "trainingFinder/internal/service/users"

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

	//authRepo := authRepository.New(pool)
	userRepo := userRepository.New(pool)
	trainingRepo := trainingRepository.New(pool)
	//authSrv := authService.New(authRepo, cfg.Server.Secret, cfg.Server.AccessTokenDuration)
	userSrv := userService.New(userRepo)
	trainingSrv := trainingService.New(trainingRepo, your_topic_name.MarshalCreateTrainingEvent)

	//ctrl := controller.New(cfg.Server, authGRPS.NewServer(authSrv), userGRPS.NewServer(userSrv), trainingGRPS.NewServer(trainingSrv))
	ctrl := controller.New(cfg.Server, userGRPS.NewServer(userSrv), trainingGRPS.NewServer(trainingSrv))
	ctrl.Run(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	log.Println("Application ex")

}
