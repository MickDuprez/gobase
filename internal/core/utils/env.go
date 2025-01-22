package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	envMap = make(map[string]string)
)

// LoadEnvFile loads environment variables from .env file
func LoadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`) // Remove quotes if present
		envMap[key] = value
	}

	return scanner.Err()
}

// GetEnvStr checks .env values first, then falls back to OS env
func GetEnvStr(key, fallback string) string {
	// Check .env map first
	if value, exists := envMap[key]; exists {
		return value
	}
	// Then check OS environment
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
