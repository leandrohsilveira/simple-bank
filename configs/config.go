package configs

import (
	"os"
)

// GetEnvOrDefault returns the value (which maybe be empty) of the environment variable named by the key or the default value if the variable is not present.
func GetEnvOrDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
