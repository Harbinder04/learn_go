package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-user-service/config"
	"go-user-service/internal/jobs"
	internal "go-user-service/internal/logger"
	"go-user-service/internal/queue"

	"github.com/redis/go-redis/v9"
)

const (
	processedKeyPrefix  = "job:processed:"
	processingKeyPrefix = "job:processing:"
	jobTTL              = 24 * time.Hour
	processingTTL       = 5 * time.Minute
)

func main() {
	cfg := config.NewConfig()

	logger := internal.NewLogger(cfg.ServerConfig.Env)
	rd, err := queue.GetRedisClient(cfg.RedisConfig, logger)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	shutSig := make(chan struct{})
	// to get a buffer time to requeue the currently running task
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Ask : Do we really need to pass shutSig as we can listen on ctx also??
	go processJobs(rd, logger, ctx)

	// Wait for shutdown signal
	go func() {
		shutDownSignal := make(chan os.Signal, 1)
		signal.Notify(shutDownSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-shutDownSignal
		logger.Info("Got the shutdown signal")
		// buffer time : Find later -> but what if the process is half done on the last second how we will roll-back??
		// 31-01-26: got answer we will hadle it using context
		// we will requeue the task when we get the shutdown signal
		cancel()
		time.Sleep(3 * time.Second)

		close(shutSig)
	}()

	<-shutSig
	logger.Info("Worker shutdown complete")
}

// but handle the requeue of task on getting a shutdown signal
func processJobs(rd *redis.Client, logger *slog.Logger, ctx context.Context) {

	logger.Info("Worker started, waiting for jobs...")

	for {
		select {
		// Todo : Current behaviour is directly shutdown no requeue logic to hadle jobs gracefully
		case <-ctx.Done():
			logger.Info("Stopping job processing...")
			return
		default:
			queueName := "Welcome Email"

			jobsData, err := rd.BRPop(ctx, 1*time.Second, queueName).Result()

			fmt.Println(jobsData)

			if err != nil {
				if errors.Is(err, context.Canceled) {
					logger.Info("Context cancelled, stopping worker")
					return
				}
				if err == redis.Nil {
					continue
				}
				logger.Error(fmt.Sprintf("Error popping job: %v", err))
				continue
			}

			if len(jobsData) < 2 {
				continue
			}

			// Parse individual job
			var job jobs.Job
			if err := json.Unmarshal([]byte(jobsData[1]), &job); err != nil {
				logger.Error("Failed to unmarshal job", "error", err, "data", jobsData[1])
				continue
			}

			logger.Info("Received job", "jobID", job.Id, "data", job.Data)

			// cehck if not job is already processed.
			processedKey := processedKeyPrefix + job.Id
			alreadyProcessed, err := rd.Exists(ctx, processedKey).Result()

			if alreadyProcessed > 0 {
				logger.Info("Job is already Processed, skipping", "jobId", job.Id)
				continue
			}

			processingKey := processingKeyPrefix + job.Id
			// SetNx : Set IF Not Exist `https://redis.io/docs/latest/commands/setnx/`
			// Here we get a lock 0 1 which is abstracted by go-redis as bool
			acquired, err := rd.SetNX(ctx, processingKey, "1", processingTTL).Result()
			if err != nil {
				logger.Error("Failed to acquire processing lock", "jobId", job.Id, "error", err)
				requeueJob(ctx, rd, queueName, job, logger)
			}

			if !acquired {
				logger.Info("Job already being processed, skipping", "jobID", job.Id)
				continue
			}

			if err := sendEmail(ctx, job, logger); err != nil {
				logger.Error("Failed to process job", "jobID", job.Id, "error", err)

				//Ask: what wil happen here if the ctx is canceled mid way??
				rd.Del(ctx, processingKey)

				requeueCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
				defer cancel()

				requeueJob(requeueCtx, rd, queueName, job, logger)
				continue
			}

			if err := rd.Set(ctx, processedKey, "1", jobTTL).Err(); err != nil {
				logger.Error("Failed to mark job as processed", "jobId", job.Id, "error", err)
			}

			rd.Del(ctx, processingKey)

			logger.Info("successfully processed job", "jobId", job.Id)
		}
	}
}

func sendEmail(ctx context.Context, jobData jobs.Job, logger *slog.Logger) error {
	logger.Info("Sending welcome email to:", jobData.Id, "")

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Simulate work
		//Todo: later send email using AWS services
		time.Sleep(2 * time.Second)
	}
	return nil
}

func requeueJob(ctx context.Context, rd *redis.Client, queueName string, job jobs.Job, logger *slog.Logger) {
	encodeJob, err := json.Marshal(job)
	if err != nil {
		logger.Error("Failed to marshal job for requeue", "jobID", job.Id, "error", err)
	}

	if err = rd.LPush(ctx, queueName, encodeJob).Err(); err != nil {
		logger.Error("Failed to requeue job", "jobId", job.Id, "error", err)
	} else {
		logger.Info("Requeue job", "jobId", job.Id)
	}
}
