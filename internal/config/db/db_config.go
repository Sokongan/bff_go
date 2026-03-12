package db_config

import (
	"fmt"
	"os"
)

type DbConfig struct {
	DSN string
}

type StoreConfig struct {
	StoreAddress  string
	StorePassword string
	StoreDB       int
}

func LoadDBConfig() (*DbConfig, error) {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DSN environment variable is not set")
	}

	return &DbConfig{
		DSN: dsn,
	}, nil
}

func LoadRedisConfig() (*StoreConfig, error) {
	storeAddress := os.Getenv("REDIS_ADDR")
	if storeAddress == "" {
		return nil, fmt.Errorf("REDIS_ADDR environment variable is not set")
	}

	storePassword := os.Getenv("REDIS_PASSWORD")
	if storePassword == "" {
		return nil, fmt.Errorf("REDIS_PASSWORD environment variable is not set")
	}

	return &StoreConfig{
		StoreAddress:  storeAddress,
		StorePassword: storePassword,
		StoreDB:       0,
	}, nil
}
