package main

import (
	"flag"
	"fmt"

	frame "github.com/karlhsu/frame"
	"github.com/karlhsu/frame/internal/app"
	"github.com/karlhsu/frame/internal/module/admin"
	"github.com/karlhsu/frame/internal/module/api"
	"github.com/karlhsu/frame/internal/server"
)

func main() {
	cfgFile := flag.String("config", "", "config file path")
	flag.Parse()

	if err := app.Init(*cfgFile); err != nil {
		panic(fmt.Sprintf("init app failed: %v", err))
	}

	if err := app.AutoMigrate(); err != nil {
		panic(fmt.Sprintf("auto migrate failed: %v", err))
	}
	app.Log.Info("database migrated")

	if err := app.SeedData(); err != nil {
		app.Log.Warnf("seed data: %v", err)
	}

	router := server.NewRouter(
		frame.AdminDist,
		admin.New(),
		api.New(),
	)

	addr := fmt.Sprintf(":%d", app.Cfg.Server.Port)
	app.Log.Infof("server starting at %s", addr)

	if err := router.Run(addr); err != nil {
		panic(fmt.Sprintf("server run failed: %v", err))
	}
}
