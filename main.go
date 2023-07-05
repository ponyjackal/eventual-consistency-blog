package main

import (
	"os"
	"time"

	"github.com/ponyjackal/eventual-consistency-blog/config"
)

func main() {
	// set timezone
	loc, _ := time.LoadLocation(os.Getenv("SERVER_TIMEZONE"))
	time.Local = loc

	// load config
	if err := config.SetupConfig(); err != nil {

	}
	masterDSN, replicaDSN := config.DatabaseConfig()

}