package config

import (
	"fmt"
	"log"
	"os"
)

type ServerConfiguration struct {
	Port 					string
	Secret					string
	LimitCountPerRequest	int64
}

func ServerConfig() string {
	appServer := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))
	log.Print("Server Running at :", appServer)
	return appServer
}