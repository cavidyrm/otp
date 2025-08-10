package services

import (
	"context"
	"testing"
	"time"

	"otp/internal/config"
	"otp/internal/models"
)

// Mock repositories for testing
type mockUserRepository struct {
	users map[string]*models.User
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *mockUserRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	for _, user := range m.users {
		if user.PhoneNumber == phoneNumber {
			return user, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *models.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) List(ctx context.Context, query models.PaginationQuery) (*models.UserListResponse, error) {
	// Mock implementation
	return &models.UserListResponse{}, nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}

type mockOTPRepository struct {
	otps map[string]*models.OTP
}

func (m *mockOTPRepository) Create(ctx context.Context, otp *models.OTP) error {
	m.otps[otp.PhoneNumber] = otp
	return nil
}

func (m *mockOTPRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.OTP, error) {
	if otp, exists := m.otps[phoneNumber]; exists {
		return otp, nil
	}
	return nil, nil
}

func (m *mockOTPRepository) MarkAsUsed(ctx context.Context, phoneNumber string) error {
	if otp, exists := m.otps[phoneNumber]; exists {
		otp.Used = true
	}
	return nil
}

func (m *mockOTPRepository) DeleteExpired(ctx context.Context) error {
	return nil
}

func (m *mockOTPRepository) GetRecentOTPCount(ctx context.Context, phoneNumber string, since time.Time) (int, error) {
	count := 0
	for _, otp := range m.otps {
		if otp.PhoneNumber == phoneNumber && otp.CreatedAt.After(since) {
			count++
		}
	}
	return count, nil
}

func TestAuthService_GenerateOTP(t *testing.T) {
	// Setup
	cfg := &config.Config{
		OTP: config.OTPConfig{
			ExpiryMinutes: 2,
			Length:        6,
		},
		RateLimit: config.RateLimitConfig{
			MaxRequests:   3,
			WindowMinutes: 10,
		},
	}

	userRepo := &mockUserRepository{users: make(map[string]*models.User)}
	otpRepo := &mockOTPRepository{otps: make(map[string]*models.OTP)}
	authService := NewAuthService(userRepo, otpRepo, cfg)

	ctx := context.Background()
	phoneNumber := "+1234567890"

	// Test successful OTP generation
	response, err := authService.GenerateOTP(ctx, phoneNumber)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Error("Expected response, got nil")
	}

	if response.Message != "OTP sent successfully" {
		t.Errorf("Expected message 'OTP sent successfully', got %s", response.Message)
	}

	if response.ExpiresIn != 2 {
		t.Errorf("Expected expires in 2 minutes, got %d", response.ExpiresIn)
	}

	// Verify OTP was created in repository
	otp, err := otpRepo.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		t.Errorf("Expected no error getting OTP, got %v", err)
	}

	if otp == nil {
		t.Error("Expected OTP to be created, got nil")
	}

	if otp.PhoneNumber != phoneNumber {
		t.Errorf("Expected phone number %s, got %s", phoneNumber, otp.PhoneNumber)
	}

	if len(otp.Code) != 6 {
		t.Errorf("Expected OTP code length 6, got %d", len(otp.Code))
	}
}

func TestAuthService_VerifyOTP(t *testing.T) {
	// Setup
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret",
			ExpiryHours: 24,
		},
	}

	userRepo := &mockUserRepository{users: make(map[string]*models.User)}
	otpRepo := &mockOTPRepository{otps: make(map[string]*models.OTP)}
	authService := NewAuthService(userRepo, otpRepo, cfg)

	ctx := context.Background()
	phoneNumber := "+1234567890"
	otpCode := "123456"

	// Create a valid OTP
	otp := models.NewOTP(phoneNumber, otpCode, 2)
	otpRepo.otps[phoneNumber] = otp

	// Test successful OTP verification
	verification := models.OTPVerification{
		PhoneNumber: phoneNumber,
		Code:        otpCode,
	}

	response, err := authService.VerifyOTP(ctx, verification)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Error("Expected response, got nil")
	}

	if response.Token == "" {
		t.Error("Expected JWT token, got empty string")
	}

	if response.User.PhoneNumber != phoneNumber {
		t.Errorf("Expected phone number %s, got %s", phoneNumber, response.User.PhoneNumber)
	}

	// Verify user was created
	user, err := userRepo.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		t.Errorf("Expected no error getting user, got %v", err)
	}

	if user == nil {
		t.Error("Expected user to be created, got nil")
	}

	if user.PhoneNumber != phoneNumber {
		t.Errorf("Expected phone number %s, got %s", phoneNumber, user.PhoneNumber)
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	// Setup
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret",
			ExpiryHours: 24,
		},
	}

	userRepo := &mockUserRepository{users: make(map[string]*models.User)}
	otpRepo := &mockOTPRepository{otps: make(map[string]*models.OTP)}
	authService := NewAuthService(userRepo, otpRepo, cfg)

	// Test invalid token
	_, err := authService.ValidateToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}

	// Test empty token
	_, err = authService.ValidateToken("")
	if err == nil {
		t.Error("Expected error for empty token, got nil")
	}
}
