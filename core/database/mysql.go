package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/MickDuprez/gobase/core/utils"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func NewDBConfig() *Config {
	isDev := utils.GetEnvBool("IS_DEV", true)

	cfg := &Config{
		Host:     utils.GetEnvStr("DB_HOST", "localhost"),
		Port:     utils.GetEnvStr("DB_PORT", "3306"),
		User:     utils.GetEnvStr("DB_USER", "root"),
		Password: utils.GetEnvStr("DB_PASSWORD", "password"),
		Database: utils.GetEnvStr("DB_NAME", "gobase"),
	}

	if !isDev {
		// Validate required settings in production
		if cfg.Password == "password" {
			log.Fatal("Production database password not set")
		}
	}

	return cfg
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.User, c.Password, c.Host, c.Port, c.Database)
}

type DB struct {
	*sql.DB
}

func New(cfg *Config) (*DB, error) {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Set sensible defaults
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &DB{db}, nil
}
