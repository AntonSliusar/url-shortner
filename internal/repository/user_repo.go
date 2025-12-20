package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"url-shortner/internal/config"
	"url-shortner/internal/models"
)

var ErrNotFound = errors.New("resource not found in database")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(cfg *config.Config) *UserRepository {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.UserDatabase.Host,
		cfg.UserDatabase.Port,
		cfg.UserDatabase.User,
		cfg.UserDatabase.Password,
		cfg.UserDatabase.DBName,
		cfg.UserDatabase.SSL,)
	db, err := sql.Open("postgres", url)
	if err != nil {
		slog.Error("Failed to connect to database:","error:", err, slog.String("UserRepositort:","failde"))
	}
	err = db.Ping()
	if err != nil {
		slog.Error("Failed to ping database:","error:", err, slog.String("UserRepositort:","ping failde"))
	}
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user models.User) error {
	query := "INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, user.Username, user.Email, user.PasswordHash, user.Role)
	if err != nil {
		slog.Error("Failed to create user:", "error:", err)
		return err
	}
	return nil
}

func (r *UserRepository) CreateGoogleUser(user models.User) error {
	query := "INSERT INTO users (username, email, google_id, role, email_verified) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, user.Username, user.Email, user.GoogleID, user.Role, user.IsVerified)
	if err != nil {
		slog.Error("Failed to create (google)user:", "error:", err)
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	query := "SELECT id, username, email, password_hash, role FROM users WHERE username = $1"
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		slog.Error("Failed to get user by username:", "error:", err)
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := "SELECT id, username, email, password_hash, role FROM users WHERE email = $1"
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return models.User{}, ErrNotFound
		}
		slog.Error("Failed to get user by email:", "error:", err)
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) EmailVerifiedTrue(email string) error{
	query := "UPDATE users SET email_verified = true WHERE email = $1"
	_, err := r.db.Exec(query, email)
	if err != nil {
		slog.Error("Failde to update email_verified", "error", err)
		return err
	}
	return nil
}
