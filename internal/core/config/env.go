package config

import (
	"fmt"
	"os"
	"strconv"
)

// GetEnvStr returns environment variable or fallback if not found
func GetEnvStr(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// GetEnvBool returns environment variable as boolean or fallback if not found
func GetEnvBool(key string, fallback bool) bool {
	if strValue, exists := os.LookupEnv(key); exists {
		if value, err := strconv.ParseBool(strValue); err == nil {
			return value
		}
	}
	return fallback
}

// GetEnvInt returns environment variable as integer or fallback if not found
func GetEnvInt(key string, fallback int) int {
	if strValue, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(strValue); err == nil {
			return value
		}
	}
	return fallback
}

// RequireEnvVars checks if all required environment variables are set
func RequireEnvVars(vars ...string) error {
	missing := []string{}

	for _, v := range vars {
		if _, exists := os.LookupEnv(v); !exists {
			missing = append(missing, v)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("required environment variables not set: %v", missing)
	}
	return nil
}
