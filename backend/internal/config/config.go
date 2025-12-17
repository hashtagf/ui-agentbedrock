package config

import (
	"os"
)

type Config struct {
	Port           string
	MongoDBURI     string
	DatabaseName   string
	AgentID        string
	AgentAliasID   string
	AgentName      string // Display name for the main agent
	AWSRegion      string
	AllowedOrigins string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		MongoDBURI:     getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DatabaseName:   getEnv("DATABASE_NAME", "agentbedrock"),
		AgentID:        getEnv("AGENT_ID", ""),
		AgentAliasID:   getEnv("AGENT_ALIAS", ""),
		AgentName:      getEnv("AGENT_NAME", "Main Agent"),
		AWSRegion:      getEnv("AWS_REGION", "us-east-1"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
