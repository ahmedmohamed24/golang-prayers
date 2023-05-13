package config

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type DB_CONFIG struct {
	DB_USER     string
	DB_PASSWORD string
	DB_PORT     string
	DB_NAME     string
	DB_HOST     string
	REDIS_HOST  string
	REDIS_PORT  string
}

type APP struct {
	DB       *gorm.DB
	APP_PORT string
	REDIS_DB *redis.Client
}
