package service

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"time"
	"url-shortner/internal/repository"
)

type OTPService struct {
	otpRepo repository.OTPRepository
	emailSender EmailSender
}

func NewOTPService(otpRepo repository.OTPRepository, emailSender EmailSender) *OTPService {
	return &OTPService{
		otpRepo: otpRepo,
		emailSender: emailSender,
	}
}

func (s *OTPService) SendEmailCode(email, purpose string) error{
	code, err := generateCode()
	if err != nil {
		return err
	}
	err = s.otpRepo.SaveOTP(email, purpose, code, 10*time.Minute) // Maybe make expiration configurable
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("Your verification code for %s:", purpose)
	body := fmt.Sprintf("Code: %s", code)
	err = s.emailSender.Send(email, subject, body)
	if err != nil {
		return err
	}
	return nil
}

func (s *OTPService) VerifyEmailCode(email, purpose, inputCode string) (bool, error) {
	storedCode, err := s.otpRepo.GetOTP(email, purpose)
	if err != nil {
		return false, err
	}
	if storedCode != inputCode {
		slog.Info("Codes dont match")
		return false, nil
	}
	s.otpRepo.DeleteOTP(email, purpose)
	return true, nil
}

func generateCode() (string, error) {
	code, err := rand.Int(rand.Reader, big.NewInt(999999))
	if err != nil {
		slog.Error("Failed to generate OTP code:", "error:", err)
		return "", err
	}
	return fmt.Sprintf("%06d", code.Int64()), nil
}
