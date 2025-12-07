package handler

import (
	"net/http"
	"url-shortner/internal/service"

	"github.com/labstack/echo/v4"
)

type URLHandler struct {
	s *service.URLService
}

type urlRequest struct {
	URL   string `json:"url"`
	Alias string `json:"alias"`
}

func NewURLHandler(s *service.URLService) *URLHandler {
	return &URLHandler{s: s}
}

func (h *URLHandler) SaveURL(c echo.Context) error {
	var req urlRequest
	if err :=c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	err := h.s.SaveURL(req.URL, req.Alias)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"result": "URL saved successfully"})
}

func (h *URLHandler) GetURL(c echo.Context) error {
	alias := c.Param("alias")
	originalURL, err := h.s.GetURL(alias)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"url": originalURL})
}

func (h * URLHandler) UpdateURL(c echo.Context) error {
	var req urlRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	err := h.s.UpdateURL(req.Alias, req.URL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"result": "URL updated successfully",
		"alias": req.Alias,
		"url": req.URL,
	})
}

func (h *URLHandler) DeleteURL(c echo.Context) error {
	alias := c.Param("alias")
	err := h.s.DeleteURL(alias)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"result": "URL deleted successfully",
		"alias": alias,
	})
}	
func (h *URLHandler) GetAllURLs(c echo.Context) error {
	urls, err := h.s.GetAllURLs()	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, urls)
}