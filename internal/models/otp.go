package models

import (
	"time"
)

type OTP struct {
	ID          string    `json:"id" db:"id"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Code        string    `json:"code" db:"code"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Used        bool      `json:"used" db:"used"`
}

type OTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type OTPVerification struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Code        string `json:"code" binding:"required"`
}

type OTPResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in_minutes"`
}

func NewOTP(phoneNumber, code string, expiryMinutes int) *OTP {
	return &OTP{
		PhoneNumber: phoneNumber,
		Code:        code,
		ExpiresAt:   time.Now().Add(time.Duration(expiryMinutes) * time.Minute),
		CreatedAt:   time.Now(),
		Used:        false,
	}
}

func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

func (o *OTP) IsValid() bool {
	return !o.Used && !o.IsExpired()
}

func (o *OTP) MarkAsUsed() {
	o.Used = true
}
