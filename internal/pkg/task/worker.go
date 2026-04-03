package task

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

// HandlerFunc is the function signature for task handlers.
type HandlerFunc func(ctx context.Context, payload []byte) error

// Worker wraps asynq.Server as the task consumer.
type Worker struct {
	server   *asynq.Server
	mux      *asynq.ServeMux
	handlers map[string]HandlerFunc
}

func NewWorker(redisAddr, password string, db int, concurrency int, queues map[string]int) *Worker {
	if concurrency <= 0 {
		concurrency = 10
	}
	if len(queues) == 0 {
		queues = map[string]int{"default": 1}
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: password,
			DB:       db,
		},
		asynq.Config{
			Concurrency: concurrency,
			Queues:      queues,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				fmt.Printf("[task error] type=%s err=%v\n", task.Type(), err)
			}),
		},
	)

	return &Worker{
		server:   srv,
		mux:      asynq.NewServeMux(),
		handlers: make(map[string]HandlerFunc),
	}
}

// Handle registers a handler for the given task type.
func (w *Worker) Handle(typeName string, handler HandlerFunc) {
	w.handlers[typeName] = handler
	w.mux.HandleFunc(typeName, func(ctx context.Context, t *asynq.Task) error {
		return handler(ctx, t.Payload())
	})
}

// Start starts the worker (blocking).
func (w *Worker) Start() error {
	return w.server.Start(w.mux)
}

// Stop gracefully shuts down the worker.
func (w *Worker) Stop() {
	w.server.Stop()
	w.server.Shutdown()
}
