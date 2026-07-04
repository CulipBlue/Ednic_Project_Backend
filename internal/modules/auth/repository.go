package auth

import (
	"context"
	"database/sql"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{db: db}
}

func (r Repository) CreateUser(ctx context.Context, user User) (User, error) {
	result, err := r.db.ExecContext(ctx, `
INSERT INTO users (name, username, email, password_hash, role, status)
VALUES (?, ?, ?, ?, ?, ?)`,
		user.Name,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Status,
	)
	if err != nil {
		return User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return User{}, err
	}

	return r.FindByID(ctx, uint64(id))
}

func (r Repository) FindByEmail(ctx context.Context, email string) (User, error) {
	return r.findOne(ctx, "SELECT id, name, username, email, password_hash, role, status, bio, phone, avatar_url, created_at, updated_at FROM users WHERE email = ? LIMIT 1", email)
}

func (r Repository) FindByID(ctx context.Context, id uint64) (User, error) {
	return r.findOne(ctx, "SELECT id, name, username, email, password_hash, role, status, bio, phone, avatar_url, created_at, updated_at FROM users WHERE id = ? LIMIT 1", id)
}

func (r Repository) findOne(ctx context.Context, query string, args ...any) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
		&user.Bio,
		&user.Phone,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}
