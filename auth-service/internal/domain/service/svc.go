package service

import (
	"context"

	"auth-svc/internal/domain/model"
)

type UserService interface {
    Verify(ctx context.Context, verifyRequest model.VerifyRequest) (model.VerifyResponse, error)
    Logout(ctx context.Context, userID int64) error 
    Register(ctx context.Context, email, password string) (int64, error)
    DeactivateUser(ctx context.Context, userID int64) error
}

