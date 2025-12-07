package server

import (
	"url-shortner/internal/handler"

	"github.com/labstack/echo/v4"
)

func NewServer(urlHandler *handler.URLHandler) *echo.Echo {
	e := echo.New()

	e.POST("/url", urlHandler.SaveURL)
	e.GET("/url/:alias", urlHandler.GetURL)
	e.PUT("/url", urlHandler.UpdateURL)
	e.DELETE("/url/:alias", urlHandler.DeleteURL)
	e.GET("/urls", urlHandler.GetAllURLs)

	return e
}