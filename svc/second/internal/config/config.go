package config

import (
	"os"
)

type Config struct {
	ServiceName      string
	ServiceVersion   string
	Environment      string
	Port             string
	LogPath          string
	DDAgentHost      string
	DatabaseName     string
	DatabaseHost     string
	DatabaseUser     string
	DatabasePassword string
	DatabasePort     string
}

func Init() Config {
	return Config{
		ServiceName:      os.Getenv("SERVICE_NAME"),
		ServiceVersion:   os.Getenv("SERVICE_VERSION"),
		Environment:      os.Getenv("ENVIRONMENT"),
		Port:             os.Getenv("PORT"),
		LogPath:          os.Getenv("LOG_PATH"),
		DDAgentHost:      os.Getenv("DD_AGENT_HOST"),
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		DatabasePort:     os.Getenv("DATABASE_PORT"),
	}
}
