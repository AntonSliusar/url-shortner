package service

import (
	"errors"
	"url-shortner/internal/models"
	"url-shortner/internal/repository"
)


type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(user models.User) error {
	return s.userRepo.CreateUser(user)
}

func (s *UserService) GetUserByUsername(username string) (models.User, error) {
	return s.userRepo.GetUserByUsername(username)
}

func (s *UserService) GetUserByEmail(email string) (models.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

func(s *UserService) EmailVerifiedTrue(email string) error {
	return s.userRepo.EmailVerifiedTrue(email)
}

func (s *UserService) HandleGoogleUser(email, name, googleID string) (models.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if errors.Is(err, repository.ErrNotFound) {
		newUser := models.User{
			Email: email,
			Username: name,
			GoogleID: googleID,
			Role: "user",
			IsVerified: true,
		}
		err = s.userRepo.CreateGoogleUser(newUser)
		user = newUser
	}
	if err != nil {
		return user, err
	}
	return user, nil
}
