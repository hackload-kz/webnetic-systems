package repository

import (
	"context"
	"errors"

	"auth-svc/internal/domain/model"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrInvalidUserData   = errors.New("invalid user data")
)

type ListParams struct {
	Limit  int
	Offset int
	Active *bool // nil - all, true - only active, false - only inactive
}

type SearchParams struct {
	Query  string
	Limit  int
	Offset int
	Active *bool
}

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	// Create(ctx context.Context, user *model.User) error
	// Update(ctx context.Context, user *model.User) error
	// Delete(ctx context.Context, id int64) error
	// List(ctx context.Context, params ListParams) ([]*model.User, int64, error)
	// Search(ctx context.Context, params SearchParams) ([]*model.User, int64, error)
	// ExistsByEmail(ctx context.Context, email string) (bool, error)
}