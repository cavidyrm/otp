package services

import (
	"context"
	"errors"

	"otp/internal/models"
	"otp/internal/repository"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (*models.UserResponse, error)
	List(ctx context.Context, query models.PaginationQuery) (*models.UserListResponse, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetByID(ctx context.Context, id string) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) List(ctx context.Context, query models.PaginationQuery) (*models.UserListResponse, error) {
	// Set default values if not provided
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	return s.userRepo.List(ctx, query)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(ctx, id)
}
