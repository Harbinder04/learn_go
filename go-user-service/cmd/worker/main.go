package main

import (
	"context"
	"fmt"
	"go-user-service/config"
	internal "go-user-service/internal/logger"
	"go-user-service/internal/queue"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.NewConfig()

	logger := internal.NewLogger(cfg.ServerConfig.Env)
	rd, err := queue.GetRedisClient(cfg.RedisConfig)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	shutSig := make(chan struct{})

	go processJobs(rd, logger, shutSig)

	// Wait for shutdown signal
	go func() {
		shutDownSignal := make(chan os.Signal, 1)
		signal.Notify(shutDownSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-shutDownSignal
		logger.Info("Got the shutdown signal")

		// buffer time : Find later -> but what if the process is half done on the last second how we will roll-back??
		time.Sleep(10 * time.Second)

		close(shutSig)
	}()

	<-shutSig
	logger.Info("Worker shutdown complete")
}

func processJobs(rd *redis.Client, logger *slog.Logger, shutSig chan struct{}) {
	ctx := context.Background()

	logger.Info("Worker started, waiting for jobs...")

	for {
		select {
		case <-shutSig:
			logger.Info("Stopping job processing...")
			return
		default:
			job, err := rd.BRPop(ctx, 2*time.Second, "Welcome Email").Result()

			fmt.Println(job)
			if err != nil {
				// redis.Nil means timeout (no jobs), just continue
				if err.Error() == "redis: nil" {
					continue
				}
				logger.Error(fmt.Sprintf("Error popping job: %v", err))
				continue
			}

			if len(job) > 1 {
				logger.Info(fmt.Sprintf("Processing job: %s", job[1]))
				
				if err := sendEmail(job[1], logger); err != nil {
					logger.Error(fmt.Sprintf("Failed to send email: %v", err))
				} else {
					logger.Info(fmt.Sprintf("Successfully sent email for: %s", job[1]))
				}
			}
		}
	}
}

func sendEmail(jobData string, logger *slog.Logger) error {
	// Simulate email sending
	logger.Info(fmt.Sprintf("Sending welcome email to: %s", jobData))
	
	// Simulate work
	time.Sleep(2 * time.Second)
	
	// In real implementation:
	// - Parse jobData (likely JSON with email, name, etc.)
	// - Use email service (SendGrid, AWS SES, etc.)
	// - Return error if sending fails
	
	return nil
}

// func porcessJob(rd *redis.Client) {
// 	ctx := context.Background()
// 	job := rd.BRPop(ctx, 5000, "Welcome Email").Result()

// 	if len(job.Val()) > 0 {
// 		sendEmail(job.Val())
// 	}
// }

// func sendEmail(jobs []string) {
	
// }
