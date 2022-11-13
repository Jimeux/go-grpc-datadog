package config

import (
	"os"
)

type Config struct {
	ServiceName       string
	ServiceVersion    string
	Environment       string
	Port              string
	LogPath           string
	ServerServiceHost string
	SecondServiceName string
}

func Init() Config {
	return Config{
		ServiceName:       os.Getenv("SERVICE_NAME"),
		ServiceVersion:    os.Getenv("SERVICE_VERSION"),
		Environment:       os.Getenv("ENVIRONMENT"),
		Port:              os.Getenv("PORT"),
		LogPath:           os.Getenv("LOG_PATH"),
		ServerServiceHost: os.Getenv("SECOND_SVC_HOST"),
		SecondServiceName: os.Getenv("SECOND_SVC_NAME"),
	}
}
