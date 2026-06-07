package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"trainingFinder/internal/app/controller"
	"trainingFinder/internal/config"
	"trainingFinder/internal/kafka"
	your_topic_name "trainingFinder/internal/kafka/producer/your-topic-name"
	"trainingFinder/internal/process/outbox"

	authGRPS "trainingFinder/internal/app/auth"
	trainingGRPS "trainingFinder/internal/app/training"
	userGRPS "trainingFinder/internal/app/users"
	authRepository "trainingFinder/internal/repository/auth"
	outboxRepository "trainingFinder/internal/repository/outbox"
	trainingRepository "trainingFinder/internal/repository/training"
	userRepository "trainingFinder/internal/repository/user"
	authService "trainingFinder/internal/service/auth"
	trainingService "trainingFinder/internal/service/training"
	userService "trainingFinder/internal/service/users"

	"github.com/go-co-op/gocron/v2"
	"github.com/golangmonster/pgxtransactor"
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
	oldPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal("Unable to create connection pool:", err)
	}
	defer oldPool.Close()
	pool := pgxtransactor.New(oldPool)

	authRepo := authRepository.New(pool)
	userRepo := userRepository.New(pool)
	trainingRepo := trainingRepository.New(pool)
	outboxRepo := outboxRepository.New(pool)
	authSrv := authService.New(authRepo, cfg.Server.Secret, cfg.Server.AccessTokenDuration)
	userSrv := userService.New(userRepo)
	trainingSrv := trainingService.New(trainingRepo, your_topic_name.MarshalCreateTrainingEvent, outboxRepo)
	cfgAuth, err := config.NewAuthConfig(cfg.Server.AuthConfigPath)
	if err != nil {
		log.Fatal("failed to load auth config: %v", err)
	}
	//ctrl := controller.New(cfg.Server, authGRPS.NewServer(authSrv), userGRPS.NewServer(userSrv), trainingGRPS.NewServer(trainingSrv))
	ctrl := controller.New(cfg.Server, *cfgAuth, userGRPS.NewServer(userSrv), trainingGRPS.NewServer(trainingSrv), authGRPS.NewServer(authSrv))
	ctrl.Run(ctx)

	outboxProcess := outbox.New(outboxRepo, nil)

	s, err := gocron.NewScheduler()
	if err != nil {
		return
	}

	if cfg.TrainingOutboxProcessEnabled {
		_, err := s.NewJob(
			gocron.DurationJob(cfg.TrainingOutboxProcessDuration),
			gocron.NewTask(outboxProcess.ProduceOutboxMessages, ctx, kafka.TrainingTopic),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		)
		if err != nil {
			return
		}
	}

	s.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	log.Println("Application ex")
}
