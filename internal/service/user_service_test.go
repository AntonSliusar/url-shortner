package service_test

import (
	"testing"
	"url-shortner/internal/models"
	"url-shortner/internal/repository"
	"url-shortner/internal/service"
	"url-shortner/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleGoogleUser(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	service := service.NewUserService(mockRepo)

	email := "test@gmail.com"
	name := "Test User"
	Id := "12345"

	t.Run("user exists", func (t *testing.T) {
		existingUser := models.User{Email: email, Username: name}
		mockRepo.On("GetUserByEmail", email).Return(existingUser, nil).Once()

		user, err := service.HandleGoogleUser(email, name, Id)
		assert.NoError(t, err)
		assert.Equal(t, email, user.Email)
		mockRepo.AssertExpectations(t)
	})
	t.Run("user does not exist - create new", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", email).Return(models.User{}, repository.ErrNotFound).Once()
		mockRepo.On("CreateGoogleUser", mock.AnythingOfType("models.User")).Return(nil).Once()

		user, err := service.HandleGoogleUser(email, name, Id)
		assert.NoError(t, err)
		assert.Equal(t, "user", user.Role)
		mockRepo.AssertExpectations(t)
	})

}