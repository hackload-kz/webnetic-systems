package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/go-faster/errors"
	"github.com/golang-jwt/jwt/v4"

	"auth-svc/internal/domain/model"
	repoModel "auth-svc/internal/domain/repository"
	"auth-svc/internal/repository"
)

type AuthService struct {
	userRepo  repository.PostgresUserRepository
	JwtSecret string
}

func NewAuthService(repo repository.PostgresUserRepository, JwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  repo,
		JwtSecret: JwtSecret,
	}
}

func (s *AuthService) verifyPassword(password, hashedPassword, salt string) bool {
	// создаем хеш из пароля + соль 
	combined := password + salt
	hash := sha256.Sum256([]byte(combined))
	computedHash := hex.EncodeToString(hash[:])

	return computedHash == hashedPassword
}


func (s *AuthService) generateToken(user *model.User) (string, error) {
	claims := model.AuthClaims{
		UserID:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JwtSecret))
}

func (s *AuthService) VerifyUser(ctx context.Context, email, password string) (*model.VerifyResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repoModel.ErrUserNotFound) {
			// не возвращаем ошибку, а возвращаем ответ с сообщением
			return &model.VerifyResponse{
				Success: false,
				Message: "Invalid credentials",
			}, nil
		}
		// для других ошибок возвращаем внутреннюю ошибку
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive {
		return &model.VerifyResponse{
			Success: false,
			Message: "Account is inactive",
		}, nil
	}

	// чек срока действия (если есть)
	if user.ExpiresAt != nil && user.ExpiresAt.Before(time.Now()) {
		return &model.VerifyResponse{
			Success: false,
			Message: "Account has expired",
		}, nil
	}

	// чек пароля
	if !s.verifyPassword(password, user.Password, user.Salt) {
		return &model.VerifyResponse{
			Success: false,
			Message: "Invalid credentials",
		}, nil
	}

	// генерируем токен для gateway
	token, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.VerifyResponse{
		Success: true,
		Message: "Authentication successful",
		User: &model.UserInfo{
			ID:        user.ID, 
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
		},
		Token: token,
	}, nil
}

// доп метод для валидации JWT токена
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*model.AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		// чек метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.JwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(*model.AuthClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id int64) (*model.UserInfo, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoModel.ErrUserNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &model.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
	}, nil
}

// Метод для получения пользователя по Email
func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*model.UserInfo, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repoModel.ErrUserNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &model.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
	}, nil
}