package app

import "github.com/karlhsu/frame/internal/pkg/task"

var TaskMgr *task.Manager

func InitTask() {
	TaskMgr = task.NewManager(task.ManagerConfig{
		RedisAddr:     Cfg.Redis.Addr(),
		RedisPassword: Cfg.Redis.Password,
		RedisDB:       Cfg.Redis.DB,
		Concurrency:   Cfg.Task.Concurrency,
		Queues:        Cfg.Task.Queues,
	})
	Log.Info("task manager initialized")
}
