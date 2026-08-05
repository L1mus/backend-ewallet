package config

import (
	"fmt"
	"os"
	"strconv"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (c DBConfig) DataSourceName() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func LoadConfiguration() (DBConfig, error) {
	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return DBConfig{}, err
	}

	cfg := DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     mustGetEnv("DB_USER"),
		Password: mustGetEnv("DB_PASS"),
		DBName:   getEnv("DB_NAME", "ewalletdb"),
		SSLMode:  getEnv("DB_SSL", "disable"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return v
}
