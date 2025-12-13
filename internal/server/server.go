package server

import (
	"fmt"
	"log/slog"
	_ "url-shortner/docs"
	"url-shortner/internal/auth"
	"url-shortner/internal/config"
	"url-shortner/internal/handler"
	"url-shortner/internal/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
)

func NewServer(urlHandler *handler.URLHandler, authHandler *auth.AuthHandler, cfg *config.Config){
	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	authRoutes := e.Group("/auth")
	authRoutes.POST("/register", authHandler.RegisterUser)
	authRoutes.POST("/login", authHandler.LoginUser)

	urlRoutes := e.Group("/urls", middleware.JWTMiddleware)
	urlRoutes.POST("", urlHandler.SaveURL)
	urlRoutes.GET("/:alias", urlHandler.GetURL)
	urlRoutes.PUT("", urlHandler.UpdateURL)
	urlRoutes.GET("", urlHandler.GetAllURLs)
	urlRoutes.DELETE("/:alias", middleware.IsAdmin(urlHandler.DeleteURL))

	adress := fmt.Sprintf("%s:%s", cfg.HTTPServer.Host, cfg.HTTPServer.Port)

	slog.Error("Server starting:", e.Start(adress))
}