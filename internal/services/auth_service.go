package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"otp/internal/config"
	"otp/internal/models"
	"otp/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	GenerateOTP(ctx context.Context, phoneNumber string) (*models.OTPResponse, error)
	VerifyOTP(ctx context.Context, verification models.OTPVerification) (*models.AuthResponse, error)
	ValidateToken(tokenString string) (*models.Claims, error)
}

type authService struct {
	userRepo repository.UserRepository
	otpRepo  repository.OTPRepository
	config   *config.Config
}

func NewAuthService(userRepo repository.UserRepository, otpRepo repository.OTPRepository, config *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		otpRepo:  otpRepo,
		config:   config,
	}
}

func (s *authService) GenerateOTP(ctx context.Context, phoneNumber string) (*models.OTPResponse, error) {
	// Check rate limiting
	since := time.Now().Add(-s.config.GetRateLimitWindow())
	count, err := s.otpRepo.GetRecentOTPCount(ctx, phoneNumber, since)
	if err != nil {
		return nil, fmt.Errorf("failed to check rate limit: %w", err)
	}

	if count >= s.config.RateLimit.MaxRequests {
		return nil, errors.New("rate limit exceeded. Please try again later")
	}

	// Generate OTP code
	code, err := s.generateRandomCode(s.config.OTP.Length)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Create OTP record
	otp := models.NewOTP(phoneNumber, code, s.config.OTP.ExpiryMinutes)
	err = s.otpRepo.Create(ctx, otp)
	if err != nil {
		return nil, fmt.Errorf("failed to save OTP: %w", err)
	}

	// Print OTP to console (for development)
	fmt.Printf("OTP for %s: %s (expires in %d minutes)\n", phoneNumber, code, s.config.OTP.ExpiryMinutes)

	return &models.OTPResponse{
		Message:   "OTP sent successfully",
		ExpiresIn: s.config.OTP.ExpiryMinutes,
	}, nil
}

func (s *authService) VerifyOTP(ctx context.Context, verification models.OTPVerification) (*models.AuthResponse, error) {
	// Get the latest valid OTP for the phone number
	otp, err := s.otpRepo.GetByPhoneNumber(ctx, verification.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get OTP: %w", err)
	}

	if otp == nil {
		return nil, errors.New("invalid or expired OTP")
	}

	// Verify OTP code
	if otp.Code != verification.Code {
		return nil, errors.New("invalid OTP code")
	}

	// Check if OTP is still valid
	if !otp.IsValid() {
		return nil, errors.New("OTP has expired")
	}

	// Mark OTP as used
	err = s.otpRepo.MarkAsUsed(ctx, verification.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	// Check if user exists
	user, err := s.userRepo.GetByPhoneNumber(ctx, verification.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Create new user if doesn't exist
	if user == nil {
		user = models.NewUser(verification.PhoneNumber)
		err = s.userRepo.Create(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Update last login time
	user.UpdateLastLogin()
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Generate JWT token
	token, expiresAt, err := s.generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		Token:     token,
		User:      user.ToResponse(),
		ExpiresAt: expiresAt,
	}, nil
}

func (s *authService) ValidateToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *authService) generateRandomCode(length int) (string, error) {
	const digits = "0123456789"
	code := make([]byte, length)

	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}

	return string(code), nil
}

func (s *authService) generateJWT(user *models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.config.GetJWTExpiry())

	claims := &models.Claims{
		UserID:      user.ID,
		PhoneNumber: user.PhoneNumber,
		Exp:         expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
