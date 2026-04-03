package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/karlhsu/frame/internal/app"
	"github.com/karlhsu/frame/internal/tasks"
)

func main() {
	cfgFile := flag.String("config", "", "config file path")
	flag.Parse()

	if err := app.Init(*cfgFile); err != nil {
		panic(fmt.Sprintf("init app failed: %v", err))
	}

	mgr := app.TaskMgr

	// Register task handlers
	tasks.RegisterHandlers(mgr.Worker)
	app.Log.Info("task handlers registered")

	// Register cron jobs
	tasks.RegisterCronJobs(mgr.Scheduler)
	app.Log.Info("cron jobs registered")

	// Start scheduler in background
	go func() {
		if err := mgr.Scheduler.Start(); err != nil {
			app.Log.Errorf("scheduler start failed: %v", err)
		}
	}()
	app.Log.Info("scheduler started")

	// Start worker (blocks until shutdown signal)
	go func() {
		if err := mgr.Worker.Start(); err != nil {
			app.Log.Errorf("worker start failed: %v", err)
		}
	}()
	app.Log.Info("worker started, waiting for tasks...")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.Log.Info("shutting down worker...")
	mgr.Scheduler.Stop()
	mgr.Worker.Stop()
	mgr.Close()
	app.Log.Info("worker stopped")
}
