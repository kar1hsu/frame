package tasks

import (
	"github.com/kar1hsu/frame/internal/pkg/task"
	"github.com/kar1hsu/frame/internal/pkg/utils"
)

// RegisterHandlers registers all task handlers on the worker.
func RegisterHandlers(w *task.Worker) {
	w.Handle(TypeEmailSend, HandleEmailSend)
	w.Handle(TypeCleanup, HandleCleanup)
}

// RegisterCronJobs registers all scheduled tasks on the scheduler.
func RegisterCronJobs(s *task.Scheduler) {
	// Daily cleanup at 2:00 AM
	s.Register(task.CronTask{
		Cron:     utils.DailyAt(2, 0),
		TypeName: TypeCleanup,
		Payload:  nil,
		Queue:    "low",
	})
}
