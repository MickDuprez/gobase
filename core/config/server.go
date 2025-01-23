package config

import "github.com/MickDuprez/gobase/core/utils"

type ServerConfig struct {
	Port string
}

func NewServerConfig() *ServerConfig {
	isDev := utils.GetEnvBool("IS_DEV", true)

	cfg := &ServerConfig{
		Port: utils.GetEnvStr("PORT", ":3000"),
	}

	if !isDev {
		// Any production-specific server settings
		if cfg.Port == ":3000" {
			cfg.Port = ":8080" // Default production port
		}
	}

	return cfg
}
