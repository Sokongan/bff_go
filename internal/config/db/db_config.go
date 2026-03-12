package db_config

import (
	"os"
)

type DbConfig struct {
	DbDSN string
}

type StoreConfig struct {
	StoreAddress  string
	StorePassword string
	StoreDB       int
}

func LoadDBConfig() *DbConfig {
	return &DbConfig{
		DbDSN: os.Getenv("DSN"),
	}
}

func LoadRedisConfig() *StoreConfig {
	return &StoreConfig{
		StoreAddress:  os.Getenv("REDIS_ADDR"),
		StorePassword: os.Getenv("REDIS_PASSWORD"),
		StoreDB:       0,
	}
}
