package service_test

import (
	"errors"
	"strings"
	"testing"
	"time"
	"url-shortner/internal/service"
	"url-shortner/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendEmailCode(t *testing.T) {
	mockRepo := new(mocks.OTPRepository)
	mockEmail := new(mocks.EmailSender)
	service := service.NewOTPService(mockRepo, mockEmail)

	email := "user@gmail.com"
	purpose := "rehistration"

	mockRepo.On("SaveOTP", email, purpose, mock.Anything, 10*time.Minute).Return(nil)
	mockEmail.On("Send", email, mock.MatchedBy(func (s string) bool {
		return strings.Contains(s, purpose)
	}), mock.Anything).Return(nil)

	err := service.SendEmailCode(email, purpose)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestVerifyEmailCode(t *testing.T) {
	mockRepo := new(mocks.OTPRepository)
	service := service.NewOTPService(mockRepo, nil)

	email := "user@gmail.com"
	purpose := "registration"

	tests := []struct{
		name string
		inputCode string
		mockReturn string
		mockErr error
		expectSuccess bool
		expectErr bool
		shouldDelete bool
	}{
		{
			name: "Success case",
			inputCode: "123456",
			mockReturn: "123456",
			mockErr: nil,
			expectSuccess: true,
			expectErr: false,
			shouldDelete: true, 
		},	
		{
			name: "Wrong code",
			inputCode: "000000",
			mockReturn: "123456",
			mockErr: nil,
			expectSuccess: false,
			expectErr: false,
			shouldDelete: false, 
		},	
		{
            name:          "Repository error (not found)",
            inputCode:     "123456",
            mockReturn:    "",
            mockErr:       errors.New("not found"),
            expectSuccess: false,
            expectErr:     true,
            shouldDelete:  false,
        },
	}

	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo.On("GetOTP", email, purpose).Return(tt.mockReturn, tt.mockErr).Once()

            if tt.shouldDelete {
                mockRepo.On("DeleteOTP", email, purpose).Return(nil).Once()
            }

            result, err := service.VerifyEmailCode(email, purpose, tt.inputCode)

            if tt.expectErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
            assert.Equal(t, tt.expectSuccess, result)

            mockRepo.AssertExpectations(t)
        })
    }
}