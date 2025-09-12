package main

import (
	"context"
	"encoding/json"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
	"user-management-api/internal/config"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"
	"user-management-api/pkg/mail"
	"user-management-api/pkg/rabbitmq"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Worker struct {
	rabbitMQ    rabbitmq.RabbitMQService
	mailService mail.EmailProviderService
	cfg         *config.Config
	logger      *zerolog.Logger
}

func NewWorker(cfg *config.Config) *Worker {
	log := utils.NewLoggerWithPath("worker.log", "info")

	// Connect RabbitMQ
	rabbitMQ, err := rabbitmq.NewRabbitMQService(
		utils.GetEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"), log,
	)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to init RabbitMQ service")
	}

	// Init MailService
	mailLogger := utils.NewLoggerWithPath("mail.log", "info")
	factory, err := mail.NewProviderFactory(mail.ProviderMailtrap)
	if err != nil {
		mailLogger.Error().Err(err).Msg("Failed to create mail provider factory")
		return nil
	}

	mailService, err := mail.NewMailService(cfg, mailLogger, factory)
	if err != nil {
		mailLogger.Error().Err(err).Msg("Failed to initialize mail service")
		return nil
	}

	return &Worker{
		rabbitMQ:    rabbitMQ,
		mailService: mailService,
		cfg:         cfg,
		logger:      log,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	const emailQueueName = "auth_email_queue"

	handler := func(body []byte) error {
		w.logger.Debug().Msgf("Received message: %s", string(body))

		var email mail.Email
		if err := json.Unmarshal(body, &email); err != nil {
			w.logger.Error().Err(err).Msg("Failed to unmarshal message")
			return err
		}

		if err := w.mailService.SendMail(ctx, &email); err != nil {
			return utils.NewError("Failed to send passwod reset email", utils.ErrCodeInternal)
		}

		w.logger.Info().Msgf("Email sent successfully to %v", email.To)

		return nil
	}

	if err := w.rabbitMQ.Consume(ctx, emailQueueName, handler); err != nil {
		w.logger.Error().Err(err).Msg("Failed to satrt consumer")
		return err
	}

	w.logger.Info().Msgf("Worker started, consuming from queue: %s", emailQueueName)

	<-ctx.Done()
	w.logger.Info().Msgf("Worker stopped consuming due to context cancellation")

	return ctx.Err()
}

func (w *Worker) Shutdown(ctx context.Context) error {
	w.logger.Info().Msg("Shutting down worker ... ")

	if err := w.rabbitMQ.Close(); err != nil {
		w.logger.Error().Err(err).Msg("Failed to close RabbitMQ")
		return err
	}

	w.logger.Info().Msg("RabbitMQ connection closed successfully")

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			w.logger.Warn().Msg("Shutdown timeout exceeded")
			return ctx.Err()
		}
	default:
	}

	w.logger.Info().Msg("Worker shutdown completed")

	return nil
}

func main() {
	rootDir := utils.MustGetWorkingDir()

	logFile := filepath.Join(rootDir, "internal/logs/app.log")

	logger.InitLogger(logger.LoggerConfig{
		Level:      "info",
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      utils.GetEnv("APP_EVN", "development"),
	})

	if err := godotenv.Load(filepath.Join(rootDir, ".env")); err != nil {
		logger.Log.Warn().Msg("⚠️ No .env file found")
	} else {
		logger.Log.Info().Msg("✅ Loaded successfully .env in worker proccess")
	}

	// Initialize configuration
	cfg := config.NewConfig()

	worker := NewWorker(cfg)
	if worker == nil {
		logger.Log.Fatal().Msg("Failed to create worker")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := worker.Start(ctx); err != nil && err != context.Canceled {
			logger.Log.Error().Err(err).Msg("Worker failed to start")
		}
	}()

	<- ctx.Done()
	logger.Log.Info().Msg("Received shutdown signal")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	if err := worker.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error().Err(err).Msg("Shutdown failed")
	}

	wg.Wait()
	logger.Log.Info().Msg("Main process terminated")
}
