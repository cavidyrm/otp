package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthResponse struct {
	Token     string       `json:"token"`
	User      UserResponse `json:"user"`
	ExpiresAt time.Time    `json:"expires_at"`
}

type Claims struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Exp         int64  `json:"exp"`
}

// GetExpirationTime implements jwt.Claims
func (c *Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.Exp, 0)), nil
}

// GetNotBefore implements jwt.Claims
func (c *Claims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

// GetIssuedAt implements jwt.Claims
func (c *Claims) GetIssuedAt() (*jwt.NumericDate, error) {
	return nil, nil
}

// GetIssuer implements jwt.Claims
func (c *Claims) GetIssuer() (string, error) {
	return "", nil
}

// GetSubject implements jwt.Claims
func (c *Claims) GetSubject() (string, error) {
	return "", nil
}

// GetAudience implements jwt.Claims
func (c *Claims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

type PaginationQuery struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Search   string `form:"search"`
}

func (p *PaginationQuery) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationQuery) GetLimit() int {
	return p.PageSize
}
