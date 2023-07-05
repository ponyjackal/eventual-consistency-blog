package database

import (
	"log"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var (
	DB	*gorm.DB
	err	error
)

// create database connection
func DatabaseConnection(masterDSN, replicaDSN string) error {
	var db = DB

	logMode, _ := strconv.ParseBool(os.Getenv("DB_LOG_MODE"))
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	logLevel := logger.Silent
	if logMode {
		logLevel = logger.Info
	}

	db, err = gorm.Open(postgres.Open(masterDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if !debug {
		db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{
				postgres.Open(replicaDSN),
			},
			Policy: dbresolver.RandomPolicy{},
		}))
	}
	if err != nil {
		log.Fatalf("DB connection error")
		return err
	}
	DB = db
	return nil
}

// GetDB connection
func GetDB() *gorm.DB {
	return DB
}