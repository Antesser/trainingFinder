package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Antesser/trainingFinder/internal/app/controller"
	"github.com/Antesser/trainingFinder/internal/config"
	"github.com/Antesser/trainingFinder/internal/kafka"
	your_topic_name "github.com/Antesser/trainingFinder/internal/kafka/producer/your-topic-name"
	"github.com/Antesser/trainingFinder/internal/process/outbox"

	authGRPC "github.com/Antesser/trainingFinder/internal/app/auth"
	bookingGRPC "github.com/Antesser/trainingFinder/internal/app/booking"
	trainingGRPC "github.com/Antesser/trainingFinder/internal/app/training"
	userGRPC "github.com/Antesser/trainingFinder/internal/app/users"
	authRepository "github.com/Antesser/trainingFinder/internal/repository/auth"
	bookingRepository "github.com/Antesser/trainingFinder/internal/repository/booking"
	outboxRepository "github.com/Antesser/trainingFinder/internal/repository/outbox"
	trainingRepository "github.com/Antesser/trainingFinder/internal/repository/training"
	userRepository "github.com/Antesser/trainingFinder/internal/repository/user"
	authService "github.com/Antesser/trainingFinder/internal/service/auth"
	bookingService "github.com/Antesser/trainingFinder/internal/service/booking"
	trainingService "github.com/Antesser/trainingFinder/internal/service/training"
	userService "github.com/Antesser/trainingFinder/internal/service/users"

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
	bookingRepo := bookingRepository.New(pool)
	userRepo := userRepository.New(pool)
	trainingRepo := trainingRepository.New(pool)
	outboxRepo := outboxRepository.New(pool)
	authSrv := authService.New(authRepo, cfg.Server.Secret, cfg.Server.AccessTokenDuration)
	userSrv := userService.New(userRepo)
	trainingSrv := trainingService.New(trainingRepo, your_topic_name.MarshalCreateTrainingEvent, outboxRepo)
	bookingSrv := bookingService.New(bookingRepo)
	cfgAuth, err := config.NewAuthConfig(cfg.Server.AuthConfigPath)
	if err != nil {
		log.Fatal("failed to load auth config: %v", err)
	}
	//ctrl := controller.New(cfg.Server, authGRPS.NewServer(authSrv), userGRPS.NewServer(userSrv), trainingGRPS.NewServer(trainingSrv))
	ctrl := controller.New(cfg.Server, *cfgAuth, userGRPC.NewServer(userSrv), trainingGRPC.NewServer(trainingSrv), authGRPC.NewServer(authSrv), bookingGRPC.NewServer(bookingSrv))
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
