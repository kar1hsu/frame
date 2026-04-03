package task

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

// CronTask defines a scheduled task with cron expression.
type CronTask struct {
	Cron     string      // cron expression, e.g. "0 2 * * *" or "@every 5m"
	TypeName string      // task type name
	Payload  interface{} // task payload (will be JSON-marshaled)
	Queue    string      // target queue (optional, defaults to "default")
}

// Scheduler wraps asynq.Scheduler for distributed cron jobs.
// Multiple Scheduler instances can run safely — only one will enqueue each cron task.
type Scheduler struct {
	scheduler *asynq.Scheduler
}

func NewScheduler(redisAddr, password string, db int) *Scheduler {
	s := asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: password,
			DB:       db,
		},
		nil,
	)
	return &Scheduler{scheduler: s}
}

// Register adds a cron task to the scheduler.
func (s *Scheduler) Register(ct CronTask) (string, error) {
	data, err := json.Marshal(ct.Payload)
	if err != nil {
		return "", fmt.Errorf("marshal cron task payload: %w", err)
	}
	task := asynq.NewTask(ct.TypeName, data)

	var opts []asynq.Option
	if ct.Queue != "" {
		opts = append(opts, asynq.Queue(ct.Queue))
	}

	entryID, err := s.scheduler.Register(ct.Cron, task, opts...)
	if err != nil {
		return "", fmt.Errorf("register cron task [%s]: %w", ct.TypeName, err)
	}
	return entryID, nil
}

// Start starts the scheduler (blocking).
func (s *Scheduler) Start() error {
	return s.scheduler.Start()
}

// Stop shuts down the scheduler.
func (s *Scheduler) Stop() {
	s.scheduler.Shutdown()
}
