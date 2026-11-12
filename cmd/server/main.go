package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	authGRPC "github.com/Antesser/trainingFinder/internal/app/auth"
	bookingGRPC "github.com/Antesser/trainingFinder/internal/app/booking"
	"github.com/Antesser/trainingFinder/internal/app/controller"
	trainingGRPC "github.com/Antesser/trainingFinder/internal/app/training"
	userGRPC "github.com/Antesser/trainingFinder/internal/app/users"
	"github.com/Antesser/trainingFinder/internal/config"
	"github.com/Antesser/trainingFinder/internal/kafka"
	"github.com/Antesser/trainingFinder/internal/kafka/producer"
	"github.com/Antesser/trainingFinder/internal/kafka/producer/training_notifications"
	"github.com/Antesser/trainingFinder/internal/process/outbox"
	authRepository "github.com/Antesser/trainingFinder/internal/repository/postgres/auth"
	bookingRepository "github.com/Antesser/trainingFinder/internal/repository/postgres/booking"
	outboxRepository "github.com/Antesser/trainingFinder/internal/repository/postgres/outbox"
	trainingRepository "github.com/Antesser/trainingFinder/internal/repository/postgres/training"
	trainingBookingRepository "github.com/Antesser/trainingFinder/internal/repository/postgres/training-booking"
	userRepository "github.com/Antesser/trainingFinder/internal/repository/postgres/user"
	"github.com/Antesser/trainingFinder/internal/repository/redis"
	authService "github.com/Antesser/trainingFinder/internal/service/auth"
	bookingService "github.com/Antesser/trainingFinder/internal/service/booking"
	trainingService "github.com/Antesser/trainingFinder/internal/service/training"
	userService "github.com/Antesser/trainingFinder/internal/service/users"
	"github.com/Antesser/trainingFinder/internal/telegram"
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
	trainingBookingRepo := trainingBookingRepository.New(pool)
	outboxRepo := outboxRepository.New(pool)
	authSrv := authService.New(authRepo, cfg.Server.Secret, cfg.Server.AccessTokenDuration)
	userSrv := userService.New(userRepo)
	trainingSrv := trainingService.New(trainingRepo, nil, outboxRepo)
	bookingSrv := bookingService.New(bookingRepo, trainingBookingRepo, training_notifications.MarshalTrainingNotification, outboxRepo)
	cfgAuth, err := config.NewAuthConfig(cfg.Server.AuthConfigPath)
	if err != nil {
		log.Printf("failed to load auth config: %v", err)
	}
	// ctrl := controller.New(cfg.Server, authGRPS.NewServer(authSrv), userGRPS.NewServer(userSrv), trainingGRPS.NewServer(trainingSrv))
	ctrl := controller.New(cfg.Server, *cfgAuth, userGRPC.NewServer(userSrv), trainingGRPC.NewServer(trainingSrv), authGRPC.NewServer(authSrv), bookingGRPC.NewServer(bookingSrv))
	ctrl.Run(ctx)

	newCache := redis.InitRedis(ctx, cfg.Server)
	tgBot, err := telegram.New(cfg.Server.BotToken, authSrv, newCache)
	if err != nil {
		log.Print(err)
		return
	}
	if err = tgBot.Run(ctx); err != nil {
		log.Print(err)
		return
	}

	producer, err := producer.New([]string{cfg.Server.KafkaHost + cfg.Server.KafkaPort})
	if err != nil {
		log.Printf("init producer: %v", err)
		return
	}
	defer func() {
		if end := producer.Close(); end != nil {
			log.Printf("close producer: %v", end)
		}
	}()

	outboxProcess := outbox.New(outboxRepo, producer)

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Printf("new scheduler: %v", err)
		return
	}

	_, err = scheduler.NewJob(
		gocron.DurationJob(cfg.Server.OutboxDuration),
		gocron.NewTask(
			outboxProcess.ProduceOutboxMessages,
			ctx,
			kafka.TrainingNotificationsTopic,
		),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		log.Printf("new workspace booking outbox job: %v", err)
		return
	}
	scheduler.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	err = producer.Close()
	if err != nil {
		return
	}
	log.Println("Application ex")
}
