package main

import (
	"fmt"
	"os"
	"platform-exercise/internal/config"
	repository "platform-exercise/internal/infra/gorm"
	"platform-exercise/internal/models"
	"platform-exercise/internal/rest"
	"platform-exercise/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setZeroLog(logLevel string) {
	if strings.ToUpper(logLevel) == "DEBUG" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return
	}

	if strings.ToUpper(logLevel) == "INFO" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		return
	}

	if strings.ToUpper(logLevel) == "WARN" {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		return
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func main() {
	cfg, err := config.Get()
	if err != nil {
		panic(fmt.Errorf("invalid value(s) retrieved from the environment: %w", err))
	}
	setZeroLog(cfg.LoggingLevel)

	db, err := gorm.Open(postgres.Open(cfg.DbURL), &gorm.Config{})
	if err != nil {
		log.Fatal().Msgf("Unable to connect to database: %v\n", err)
	}

	// Run migration
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal().Msgf("migration failed:%v", err)
	}

	repo := repository.NewUserRepository(db)

	tokenService := services.NewAuthService(os.Getenv("JWT_KEY"))
	userService := services.NewUserService(repo, tokenService)

	router := rest.NewUserRoutes(cfg, userService)

	r := gin.Default()

	router.ImportRoutes(r)

	r.Run(":8080")
}
