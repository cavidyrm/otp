package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	OTP       OTPConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret      string
	ExpiryHours int
}

type OTPConfig struct {
	ExpiryMinutes int
	Length        int
}

type RateLimitConfig struct {
	MaxRequests   int
	WindowMinutes int
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "otp_user"),
			Password: getEnv("DB_PASSWORD", "otp_password"),
			Name:     getEnv("DB_NAME", "otp_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
			ExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		},
		OTP: OTPConfig{
			ExpiryMinutes: getEnvAsInt("OTP_EXPIRY_MINUTES", 2),
			Length:        getEnvAsInt("OTP_LENGTH", 6),
		},
		RateLimit: RateLimitConfig{
			MaxRequests:   getEnvAsInt("RATE_LIMIT_MAX_REQUESTS", 3),
			WindowMinutes: getEnvAsInt("RATE_LIMIT_WINDOW_MINUTES", 10),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (c *Config) GetDatabaseURL() string {
	return "postgres://" + c.Database.User + ":" + c.Database.Password + "@" + c.Database.Host + ":" + c.Database.Port + "/" + c.Database.Name + "?sslmode=" + c.Database.SSLMode
}

func (c *Config) GetJWTExpiry() time.Duration {
	return time.Duration(c.JWT.ExpiryHours) * time.Hour
}

func (c *Config) GetOTPExpiry() time.Duration {
	return time.Duration(c.OTP.ExpiryMinutes) * time.Minute
}

func (c *Config) GetRateLimitWindow() time.Duration {
	return time.Duration(c.RateLimit.WindowMinutes) * time.Minute
}
