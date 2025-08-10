package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          string     `json:"id" db:"id"`
	PhoneNumber string     `json:"phone_number" db:"phone_number"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

type UserCreate struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type UserResponse struct {
	ID          string     `json:"id"`
	PhoneNumber string     `json:"phone_number"`
	CreatedAt   time.Time  `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

func NewUser(phoneNumber string) *User {
	now := time.Now()
	return &User{
		ID:          uuid.New().String(),
		PhoneNumber: phoneNumber,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		PhoneNumber: u.PhoneNumber,
		CreatedAt:   u.CreatedAt,
		LastLoginAt: u.LastLoginAt,
	}
}

func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}
