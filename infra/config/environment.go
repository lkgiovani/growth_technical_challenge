package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server struct {
		Port string
		Mode string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
	}
	Cache struct {
		TTL int
	}
}

func NewConfig() (*Config, error) {
	serverPort, err := getString("APP_PORT", "8080")
	if err != nil {
		return nil, err
	}

	serverMode, err := getString("GIN_MODE", "debug")
	if err != nil {
		return nil, err
	}

	dbHost, err := getString("DB_HOST", "localhost")
	if err != nil {
		return nil, err
	}

	dbPort, err := getString("DB_PORT", "5432")
	if err != nil {
		return nil, err
	}

	dbUser, err := getString("POSTGRES_USER", "postgres")
	if err != nil {
		return nil, err
	}

	dbPassword, err := getString("POSTGRES_PASSWORD", "postgres")
	if err != nil {
		return nil, err
	}

	dbName, err := getString("POSTGRES_DB", "growth_db")
	if err != nil {
		return nil, err
	}

	dbSSLMode, err := getString("DB_SSLMODE", "disable")
	if err != nil {
		return nil, err
	}

	redisHost, err := getString("REDIS_HOST", "localhost")
	if err != nil {
		return nil, err
	}

	redisPort, err := getString("REDIS_PORT", "6379")
	if err != nil {
		return nil, err
	}

	redisPassword, err := getString("REDIS_PASSWORD", "")
	if err != nil {
		return nil, err
	}

	cacheTTL, err := getInt("CACHE_TTL", 300)
	if err != nil {
		return nil, err
	}

	return &Config{
		Server: struct {
			Port string
			Mode string
		}{
			Port: serverPort,
			Mode: serverMode,
		},
		Database: struct {
			Host     string
			Port     string
			User     string
			Password string
			DBName   string
			SSLMode  string
		}{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
			SSLMode:  dbSSLMode,
		},
		Redis: struct {
			Host     string
			Port     string
			Password string
		}{
			Host:     redisHost,
			Port:     redisPort,
			Password: redisPassword,
		},
		Cache: struct {
			TTL int
		}{
			TTL: cacheTTL,
		},
	}, nil
}

func getString(key, defaultValue string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue, nil
	}
	return value, nil
}

func getInt(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s must be a valid integer", key)
	}
	return intValue, nil
}

func getBool(key string, defaultValue bool) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("environment variable %s must be a valid boolean", key)
	}
	return boolValue, nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
