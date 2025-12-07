package main

import (
	"fmt"
	"url-shortner/internal/config"
	"url-shortner/internal/logger"
	"url-shortner/internal/repository"
)

func main() {
	logger.InitLogger("local")
	cfg := config.LoadConfig()
	dbStorage := repository.NewStorage(*cfg)
	// Example usage
	dbStorage.SaveURL("https://example.com", "exmpl")
	originalURL := dbStorage.GetURL("exmpl")
	fmt.Println("Original URL:", originalURL)

	dbStorage.UpdateURL("exmpl", "https://example.org")
	updatedURL := dbStorage.GetURL("exmpl")
	fmt.Println("Updated URL:", updatedURL)
	
	///
}