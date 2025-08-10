package handlers

import (
	"net/http"

	"otp/internal/models"
	"otp/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// GenerateOTP godoc
// @Summary Generate OTP for phone number
// @Description Generate a new OTP code for the provided phone number
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.OTPRequest true "Phone number"
// @Success 200 {object} models.OTPResponse
// @Failure 400 {object} ErrorResponse
// @Failure 429 {object} ErrorResponse
// @Router /auth/otp/generate [post]
func (h *AuthHandler) GenerateOTP(c *gin.Context) {
	var request models.OTPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	response, err := h.authService.GenerateOTP(c.Request.Context(), request.PhoneNumber)
	if err != nil {
		if err.Error() == "rate limit exceeded. Please try again later" {
			c.JSON(http.StatusTooManyRequests, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate OTP"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// VerifyOTP godoc
// @Summary Verify OTP and authenticate user
// @Description Verify OTP code and authenticate/register user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.OTPVerification true "OTP verification"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/otp/verify [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var request models.OTPVerification
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	response, err := h.authService.VerifyOTP(c.Request.Context(), request)
	if err != nil {
		if err.Error() == "invalid or expired OTP" ||
			err.Error() == "invalid OTP code" ||
			err.Error() == "OTP has expired" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to verify OTP"})
		return
	}

	c.JSON(http.StatusOK, response)
}

type ErrorResponse struct {
	Error string `json:"error"`
}
