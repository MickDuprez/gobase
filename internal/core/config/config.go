package config

import (
	"github.com/MickDuprez/gobase/internal/core/database"
	"github.com/MickDuprez/gobase/internal/core/middleware"
	"github.com/MickDuprez/gobase/internal/core/utils"
)

type AppConfig struct {
	Server    *ServerConfig
	DBConfig  *database.Config
	SecConfig *middleware.SecurityConfig
}

func NewAppConfig() *AppConfig {
	isDev := utils.GetEnvBool("IS_DEV", true)

	if isDev {
		return &AppConfig{
			Server:    NewServerConfig(),
			DBConfig:  database.NewDBConfig(),
			SecConfig: middleware.NewDevSecurityConfig(),
		}
	}

	return &AppConfig{
		Server:    NewServerConfig(),
		DBConfig:  database.NewDBConfig(),
		SecConfig: middleware.NewProdSecurityConfig(),
	}
}
