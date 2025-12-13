package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get("role").(string)
		if role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]string{"erroe": "Admin access required"})
		}	
		return next(c)
	}

}