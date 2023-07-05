package main

import (
	"os"
	"time"

	"github.com/ponyjackal/eventual-consistency-blog/config"
	"github.com/ponyjackal/eventual-consistency-blog/infra/database"
	"github.com/ponyjackal/eventual-consistency-blog/infra/logger"
	"github.com/ponyjackal/eventual-consistency-blog/migrations"
	"github.com/ponyjackal/eventual-consistency-blog/routers"
	"github.com/ponyjackal/eventual-consistency-blog/services"
)

func main() {
	// set timezone
	loc, _ := time.LoadLocation(os.Getenv("SERVER_TIMEZONE"))
	time.Local = loc

	// load config
	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}

	masterDSN, replicaDSN := config.DatabaseConfig()
	if err := database.DatabaseConnection(masterDSN, replicaDSN); err != nil {
		logger.Fatalf("database DatabaseConnection error: %s", err)
	}

	migrations.Migrate()

	cm, _ := services.NewCacheManager()
	// defer closeCacheManager()
	p, _ := services.NewPublisher()
	// defer closePublisher()

	// run background services
	go cm.Run()
	go p.Run()

	router := routers.SetupRoute()
	logger.Fatalf("%v", router.Run(config.ServerConfig()))
}