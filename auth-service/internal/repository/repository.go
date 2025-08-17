package repository

import (
	"context"
	"database/sql"

	"github.com/go-faster/errors"

	"auth-svc/internal/domain/model"
	"auth-svc/internal/domain/repository"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

const (
	baseSelectQuery = `SELECT id, email, password_hash, salt, first_name, last_name, 
		birth_date, created_at, is_active, expires_at FROM users`
)

func (r *PostgresUserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := baseSelectQuery + ` WHERE id = $1`
	
	user := &model.User{}
	var birthDate, expiresAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Salt,
		&user.FirstName,
		&user.LastName,
		&birthDate,
		&user.CreatedAt,
		&user.IsActive,
		&expiresAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, errors.Wrap(err, "failed to get user by id")
	}
	
	if birthDate.Valid {
		user.BirthDate = &birthDate.Time
	}
	if expiresAt.Valid {
		user.ExpiresAt = &expiresAt.Time
	}
	
	return user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := baseSelectQuery + ` WHERE email = $1`
	
	user := &model.User{}
	var birthDate, expiresAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Salt,
		&user.FirstName,
		&user.LastName,
		&birthDate,
		&user.CreatedAt,
		&user.IsActive,
		&expiresAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, errors.Wrap(err, "failed to get user by email")
	}
	
	if birthDate.Valid {
		user.BirthDate = &birthDate.Time
	}
	if expiresAt.Valid {
		user.ExpiresAt = &expiresAt.Time
	}
	
	return user, nil
}

