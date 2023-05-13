package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ahmedmohamed24/azan/config"
	"github.com/ahmedmohamed24/azan/models"
	"github.com/ahmedmohamed24/azan/repository"
	"github.com/ahmedmohamed24/azan/routes"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	dbConfig := config.DB_CONFIG{
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_HOST:     os.Getenv("DB_HOST"),
		REDIS_HOST:  os.Getenv("REDIS_HOST"),
		REDIS_PORT:  os.Getenv("REDIS_PORT"),
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Cairo", dbConfig.DB_HOST, dbConfig.DB_USER, dbConfig.DB_PASSWORD, dbConfig.DB_NAME, dbConfig.DB_PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.CountryMethod{})
	db.AutoMigrate(&models.Seeding{})
	repository.CountryMethodSeeding(db)

	redisDB := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", dbConfig.REDIS_HOST, dbConfig.REDIS_PORT),
		Password: "",
		DB:       0,
	})

	app := config.APP{
		DB:       db,
		REDIS_DB: redisDB,
		APP_PORT: os.Getenv("APP_PORT"),
	}

	engine := gin.Default()
	v1 := engine.Group("api/v1")
	{
		routes.API_V1(v1, &app)
	}

	log.Fatal(engine.Run(app.APP_PORT))
}
