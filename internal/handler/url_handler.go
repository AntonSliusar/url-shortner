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

// @Summary SaveURL
// @Description Save a new URL with an alias
// @Tags URL
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param urlRequest body urlRequest true "URL Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /urls [post]
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

// @Summary GetURL
// @Description Send alias to get the original URL
// @Tags URL
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param alias path string true "Alias"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /urls/{alias} [get]
func (h *URLHandler) GetURL(c echo.Context) error {
	alias := c.Param("alias")
	originalURL, err := h.s.GetURL(alias)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"url": originalURL})
}

// @Summary UpdateURL
// @Description Update the original URL for a given alias
// @Tags URL
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param urlRequest body urlRequest true "URL Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /urls [put]
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

// @Summary DeleteURL
// @Description Delete a URL by its alias, only admin can delete urls
// @Tags URL
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param alias path string true "Alias"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /urls/{alias} [delete]
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

// @Summary GetAllURLs
// @Description Retrieve all stored URLs with their aliases
// @Tags URL
// @Security ApiKeyAuth
// @Accept json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /urls [get]
func (h *URLHandler) GetAllURLs(c echo.Context) error {
	urls, err := h.s.GetAllURLs()	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, urls)
}