package main

import "os"

type Config struct {
	Prefix     string
	SessionDB  string
	LogLevel   string
	BotName    string
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func LoadConfig() Config {
	return Config{
		Prefix:    getenv("BOT_PREFIX", "!"),
		SessionDB: getenv("BOT_SESSION_DB", "session.db"),
		LogLevel:  getenv("BOT_LOG_LEVEL", "INFO"),
		BotName:   getenv("BOT_NAME", "WaBot"),
	}
}
