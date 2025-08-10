package repository

import (
	"context"
	"database/sql"
	"fmt"

	"otp/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	List(ctx context.Context, query models.PaginationQuery) (*models.UserListResponse, error)
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, phone_number, created_at, updated_at, last_login_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.PhoneNumber, user.CreatedAt, user.UpdatedAt, user.LastLoginAt)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, phone_number, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	query := `
		SELECT id, phone_number, created_at, updated_at, last_login_at
		FROM users
		WHERE phone_number = $1
	`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, phoneNumber).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET phone_number = $2, updated_at = $3, last_login_at = $4
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.PhoneNumber, user.UpdatedAt, user.LastLoginAt)
	return err
}

func (r *userRepository) List(ctx context.Context, query models.PaginationQuery) (*models.UserListResponse, error) {
	// Build the base query
	baseQuery := "FROM users"
	whereClause := ""
	args := []interface{}{}

	// Add search condition if provided
	if query.Search != "" {
		whereClause = "WHERE phone_number ILIKE $1"
		args = append(args, "%"+query.Search+"%")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) %s %s", baseQuery, whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Calculate pagination
	offset := query.GetOffset()
	limit := query.GetLimit()
	totalPages := (total + limit - 1) / limit

	// Build the main query
	mainQuery := fmt.Sprintf(`
		SELECT id, phone_number, created_at, updated_at, last_login_at
		%s %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, baseQuery, whereClause, len(args)+1, len(args)+2)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, mainQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.PhoneNumber,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.LastLoginAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user.ToResponse())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &models.UserListResponse{
		Users:      users,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
