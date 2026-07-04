package account

import (
	"context"
	"database/sql"

	"github.com/CulipBlue/backend_ednic/internal/modules/auth"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{db: db}
}

func (r Repository) FindUserByID(ctx context.Context, userID uint64) (auth.User, error) {
	authRepo := auth.NewRepository(r.db)
	return authRepo.FindByID(ctx, userID)
}

func (r Repository) UpdateProfile(ctx context.Context, userID uint64, request UpdateProfileRequest) (auth.User, error) {
	_, err := r.db.ExecContext(ctx, `
UPDATE users
SET name = ?, username = ?, bio = ?, phone = ?, avatar_url = ?
WHERE id = ?`,
		request.Name,
		request.Username,
		request.Bio,
		request.Phone,
		request.AvatarURL,
		userID,
	)
	if err != nil {
		return auth.User{}, err
	}

	return r.FindUserByID(ctx, userID)
}

func (r Repository) UpdatePasswordHash(ctx context.Context, userID uint64, passwordHash string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET password_hash = ? WHERE id = ?", passwordHash, userID)
	return err
}
