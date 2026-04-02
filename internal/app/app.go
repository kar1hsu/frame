package app

import "fmt"

func Init(cfgFile string) error {
	if err := InitConfig(cfgFile); err != nil {
		return fmt.Errorf("init config: %w", err)
	}

	InitLogger()
	Log.Info("config loaded successfully")

	if err := InitDatabase(); err != nil {
		return fmt.Errorf("init database: %w", err)
	}
	Log.Info("database connected")

	if err := InitRedis(); err != nil {
		return fmt.Errorf("init redis: %w", err)
	}
	Log.Info("redis connected")

	if err := InitCasbin(); err != nil {
		return fmt.Errorf("init casbin: %w", err)
	}
	Log.Info("casbin initialized")

	return nil
}
