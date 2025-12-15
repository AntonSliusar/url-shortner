package main

import (
	_ "url-shortner/docs"
	"url-shortner/internal/auth"
	"url-shortner/internal/config"
	"url-shortner/internal/handler"
	"url-shortner/internal/logger"
	"url-shortner/internal/repository"
	"url-shortner/internal/server"
	"url-shortner/internal/service"
)

// @title URL-Shortner API
// @version 1.0
// @description This is a simple URL shortener service API.

// @host localhost:8080
// @BasePath /

// @SecurityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logger.InitLogger("local")
	cfg := config.LoadConfig()

	urlRepo := repository.NewURLRepository(cfg)
	urlService := service.NewService(urlRepo)
	urlHandler := handler.NewURLHandler(urlService)

	userRepo := repository.NewUserRepository(cfg)
	userService := service.NewUserService(userRepo)
	otpRepo := repository.NewOTPRepository()
	emailSender := service.NewSMPTSender()
	otpService := service.NewOTPService(otpRepo, emailSender)
	authHandler := auth.NewUserHandler(userService, otpService)
	server.NewServer(urlHandler, authHandler, cfg)	
}
