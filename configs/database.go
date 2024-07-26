package configs

import (
	"fmt"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SSLMode  string
	MaxConns int
}

func (dbConfig DatabaseConfig) GetDbSource() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DbName,
		dbConfig.SSLMode,
	)
}

func NewDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     GetEnvOrDefault("DATABASE_HOST", "localhost"),
		Port:     GetEnvOrDefault("DATABASE_PORT", "5432"),
		User:     GetEnvOrDefault("DATABASE_USER", "app"),
		Password: GetEnvOrDefault("DATABASE_PASSWORD", "password"),
		DbName:   GetEnvOrDefault("DATABASE_DBNAME", "app"),
		SSLMode:  GetEnvOrDefault("DATABASE_SSLMODE", "disable"),
		MaxConns: 30,
	}
}

func NewTestingDatabaseConfig() DatabaseConfig {
	config := NewDatabaseConfig()
	return DatabaseConfig{
		Host:     GetEnvOrDefault("TESTING_DATABASE_HOST", config.Host),
		Port:     GetEnvOrDefault("TESTING_DATABASE_PORT", config.Port),
		User:     GetEnvOrDefault("TESTING_DATABASE_USER", config.User),
		Password: GetEnvOrDefault("TESTING_DATABASE_PASSWORD", config.Password),
		DbName:   GetEnvOrDefault("TESTING_DATABASE_DBNAME", "app_test"),
		SSLMode:  GetEnvOrDefault("TESTING_DATABASE_SSLMODE", config.SSLMode),
		MaxConns: 15,
	}
}