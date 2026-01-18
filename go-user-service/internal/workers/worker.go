package workers

import (
	"go-user-service/internal/jobs"
	"log/slog"
	"sync"
	"time"
)

type Worker struct {
	jobQueue chan jobs.Job
	logger *slog.Logger
	wg sync.WaitGroup
}

func NewWorker(jobQueue chan jobs.Job, logger *slog.Logger) *Worker {
	return &Worker{
		jobQueue: jobQueue,
		logger: logger,
	}
}

func (w *Worker) Start() {
	w.logger.Info("Worker started, waiting for jobs...")

	for job := range w.jobQueue {
		 w.wg.Add(1)
		w.processJob(job)
	}

	w.logger.Info("Worker stopped (channel closed)")
}


func (w *Worker) processJob(job jobs.Job) {
	defer w.wg.Done()

	switch job.Type {
	case "Welcome Email":
		w.handleWelcomeEmail(job)
	case "Audit":
		w.handleAuditLog(job)
	default:
		w.logger.Warn("Unknown job type", "type", job.Type)
	}
}

func (w *Worker) handleWelcomeEmail(job jobs.Job) {
	email := job.PayLoad
	
	w.logger.Info("sending welcome email", "email", email)
	time.Sleep(2 * time.Second)

	w.logger.Info("Welcome email sent to", "email", email)
}

func (w *Worker) handleAuditLog(job jobs.Job) {
	w.logger.Info("Logging audit event", "payload", job.PayLoad)
    time.Sleep(500 * time.Millisecond)
}

func (w *Worker) Shutdown() {
	w.logger.Info("Waiting for worker to finishin-progress jobs...")
    w.wg.Wait()
    w.logger.Info("All jobs completed")
}