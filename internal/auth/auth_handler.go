package auth

import (
	"net/http"
	"url-shortner/internal/models"
	"url-shortner/internal/service"
	"url-shortner/internal/utils"

	"github.com/labstack/echo/v4"
)
const (
	registPurpose = "registration"
	loginPurpose = "logein"
)
	

type AuthHandler struct {
	userService *service.UserService
	otpService *service.OTPService
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserHandler(userService *service.UserService, otpService *service.OTPService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		otpService: otpService,
	}
}

// @Summary RegisterUser
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "user registration request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) RegisterUser(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	hashPass, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashPass,
		Role:         req.Role,
	}
	err = h.userService.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	err = h.otpService.SendEmailCode(user.Email, registPurpose)
	if err != nil {	
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"result": "Check email code to complete registration"})
}

// @Summary ConfirmRegistration
// @Description Confirm registration with email code
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.VerifyCodeInput true "email and code"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/verify-registration [post]
func (h *AuthHandler) VerifyRegistration(c echo.Context) error {
	var input models.VerifyCodeInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string] string {"error": err.Error()})
	}
	res, err := h.otpService.VerifyEmailCode(input.Email, registPurpose, input.Code)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error" : err.Error()})
	}
	if !res {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error" : "invalid code"})
	}
	h.userService.EmailVerifiedTrue(input.Email)
	return c.JSON(http.StatusOK, map[string]string{"result" : "Registration complete"})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "login credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) LoginUser(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}
	err = h.otpService.SendEmailCode(user.Email, registPurpose)
	if err != nil {	
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"result": "Check email code to complete login"})
}

// @Summary ConfirmLogin
// @Description Confirm login with email code
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.VerifyCodeInput true "email and code"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/verify-login [post]
func (h *AuthHandler) VerifyLogin(c echo.Context) error {
	var input models.VerifyCodeInput 
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string {"error": err.Error()})
	}
	user, err := h.userService.GetUserByEmail(input.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email"})
	}
	res, err := h.otpService.VerifyEmailCode(input.Email, loginPurpose, input.Code)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	if !res {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	token, err := GenerateToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}