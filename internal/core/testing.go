package core

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
	"github.com/yaza-putu/golang-starter-mongo-api/internal/database"
	_ "github.com/yaza-putu/golang-starter-mongo-api/internal/database/migrations"
)

func EnvTesting() error {
	_, b, _, _ := runtime.Caller(0)

	// Root folder of this project
	Root := filepath.Join(filepath.Dir(b), "../..")
	viper.SetConfigName(".env.test")
	viper.SetConfigType("env")
	viper.AddConfigPath(Root)
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	viper.Set("app_debug", false)
	viper.Set("app_status", "test")
	// force auto migrate true
	viper.Set("db_auto_migrate", true)

	// call database
	Mongo()

	// run server
	go HttpServe()

	return err
}

func EnvRollback() {
	database.DownMigration()
}
