package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"url-shortner/internal/config"

	_ "github.com/lib/pq"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(cfg *config.Config) *URLRepository {
	url := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.URLDatabase.Host,
		cfg.URLDatabase.Port,
		cfg.URLDatabase.User,
		cfg.URLDatabase.Password,
		cfg.URLDatabase.DBName,
		cfg.URLDatabase.SSL,
	)
	db, err := sql.Open("postgres", url)
	if err != nil {
		slog.Error("Failed to connect to database:","error:", err)
	}
	err = db.Ping()
	if err != nil {
		slog.Error("Failed to ping database:","error:", err)
	}
	return &URLRepository{db: db}
}

func (r *URLRepository) SaveURL(urlToSave string, alias string) error {
	query := "INSERT INTO urls (original_url, alias) VALUES ($1, $2)"
	_, err := r.db.Exec(query, urlToSave, alias)
	if err != nil {
		slog.Error("Failed to save URL:", "error:", err)
		return err
	}
	return nil
}

func (r *URLRepository) GetURL(alias string) (string, error) {
	var originalURL string
	query := "SELECT original_url FROM urls WHERE alias = $1"
	err := r.db.QueryRow(query, alias).Scan(&originalURL)
	if err != nil {
		slog.Error("Failed to get URL:", "error:", err)
		return "", err
	}
	return originalURL, nil
}

func (r *URLRepository) UpdateURL(alias string, newURL string) error{
	query := "UPDATE urls SET original_url = $1 WHERE alias = $2"
	_, err := r.db.Exec(query, newURL, alias)
	if err != nil {
		slog.Error("Failed to update URL:", "error:", err)
		return err
	}
	return nil
}	

func (r *URLRepository) DeleteURL(alias string) error{
	query := "DELETE FROM urls WHERE alias = $1"
	_, err := r.db.Exec(query, alias)
	if err != nil {
		slog.Error("Failed to delete URL:", "error:", err)
		return err
	}
	return nil
}

func (r *URLRepository) GetAllURLs() (map[string]string, error) {
	urls := make(map[string]string)
	query := "SELECT alias, original_url FROM urls"
	rows, err := r.db.Query(query)
	if err != nil {
		slog.Error("Failed to get all URLs:", "error:", err)
		return urls, err
	}
	defer rows.Close()

	for rows.Next() {
		var alias, originalURL string
		err := rows.Scan(&alias, &originalURL)
		if err != nil {
			slog.Error("Failed to scan row:", "error:", err)
			continue
		}
		urls[alias] = originalURL
	}
	return urls, nil
}

