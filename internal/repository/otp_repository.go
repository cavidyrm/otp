package repository

import (
	"context"
	"database/sql"
	"time"

	"otp/internal/models"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *models.OTP) error
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.OTP, error)
	MarkAsUsed(ctx context.Context, phoneNumber string) error
	DeleteExpired(ctx context.Context) error
	GetRecentOTPCount(ctx context.Context, phoneNumber string, since time.Time) (int, error)
}

type otpRepository struct {
	db *sql.DB
}

func NewOTPRepository(db *sql.DB) OTPRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) Create(ctx context.Context, otp *models.OTP) error {
	query := `
		INSERT INTO otps (phone_number, code, expires_at, created_at, used)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, otp.PhoneNumber, otp.Code, otp.ExpiresAt, otp.CreatedAt, otp.Used)
	return err
}

func (r *otpRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.OTP, error) {
	query := `
		SELECT phone_number, code, expires_at, created_at, used
		FROM otps
		WHERE phone_number = $1 AND used = false AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`
	otp := &models.OTP{}
	err := r.db.QueryRowContext(ctx, query, phoneNumber).Scan(
		&otp.PhoneNumber,
		&otp.Code,
		&otp.ExpiresAt,
		&otp.CreatedAt,
		&otp.Used,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return otp, nil
}

func (r *otpRepository) MarkAsUsed(ctx context.Context, phoneNumber string) error {
	query := `
		UPDATE otps
		SET used = true
		WHERE phone_number = $1 AND used = false
	`
	_, err := r.db.ExecContext(ctx, query, phoneNumber)
	return err
}

func (r *otpRepository) DeleteExpired(ctx context.Context) error {
	query := `
		DELETE FROM otps
		WHERE expires_at < NOW()
	`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *otpRepository) GetRecentOTPCount(ctx context.Context, phoneNumber string, since time.Time) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM otps
		WHERE phone_number = $1 AND created_at >= $2
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, phoneNumber, since).Scan(&count)
	return count, err
}
