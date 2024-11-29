package config

import (
	"os"
	"strings"
)

const (
	osPrefix = "DISCORD"
)

type Config struct {
	ChannelPrefix    string
	ConnectionString string
	Username         string
	Password         string
	LogLevel         string
	LogJSON          bool
	DiscordToken     string
}

func Load() (Config, error) {
	return Config{
		ChannelPrefix:    getEnv("BROKER_SENDING_PREFIX", "iggy.discord"),
		ConnectionString: getEnv("BROKER_CONNECTION_STRING", "nats://0.0.0.0:4222"),
		Username:         getEnv("BROKER_USERNAME", ""),
		Password:         getEnv("BROKER_PASSWORD", ""),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		LogJSON:          toBool(getEnv("LOG_JSON", "false")),
		DiscordToken:     getEnv("TOKEN", ""),
	}, nil
}

func getEnv(key string, defaultValue string) string {
	fullKey := osPrefix + "_" + key

	val := os.Getenv(fullKey)
	if val == "" {
		if defaultValue != "" {
			return defaultValue
		}
	}

	return val
}

func toBool(s string) bool {
	return strings.ToLower(s) == "true"
}
