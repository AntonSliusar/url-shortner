package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPRepository interface {
	SaveOTP(email, purpose, code string, expiration time.Duration) error
	GetOTP(email, purpose string) (string, error)
	DeleteOTP(email, purpose string) error
}

type OTPRepo struct {
	rdb *redis.Client
}
//TODO: add config to initialize redis client
func NewOTPRepository() *OTPRepo { 
	rbd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &OTPRepo{
		rdb: rbd,
	}
}

func (r *OTPRepo) SaveOTP(email, purpose, code string, expiration time.Duration) error {
	key := buildKey(email, purpose)
	err := r.rdb.Set(context.Background(), key, code, expiration).Err()
	if err != nil {
		slog.Error("Failed to save OTP:", "error:", err)
		return err
	}
	return nil
}

func (r *OTPRepo) GetOTP(email, purpose string) (string, error) {
	key := buildKey(email, purpose)
	res, err := r.rdb.Get(context.Background(), key).Result()
	if err == redis.Nil {
		slog.Error("OTP not found for key:", "key:", key)
		return "", err
	}
	if err != nil {
		slog.Error("Failed to get OTP:", "error:", err)
		return "", err
	
	}
	return res, nil
}

func (r *OTPRepo) DeleteOTP(email, purpose string) error {
	key := buildKey(email, purpose)
	err := r.rdb.Del(context.Background(), key).Err()
	if err != nil {
		slog.Error("Failed to delete OTP:", "error:", err)
		return err
	}
	return nil
}

func buildKey(email, purpose string) string {
	return email + ":" + purpose
}