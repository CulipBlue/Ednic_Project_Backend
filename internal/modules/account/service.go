package account

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCurrentPassword = errors.New("invalid current password")

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) GetProfile(ctx context.Context, userID uint64) (ProfileResponse, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return ProfileResponse{}, err
	}

	return ProfileResponse{User: user}, nil
}

func (s Service) UpdateProfile(ctx context.Context, userID uint64, request UpdateProfileRequest) (ProfileResponse, error) {
	request.Name = strings.TrimSpace(request.Name)
	request.Username = strings.TrimSpace(request.Username)
	request.Bio = trimOptional(request.Bio)
	request.Phone = trimOptional(request.Phone)
	request.AvatarURL = trimOptional(request.AvatarURL)

	user, err := s.repo.UpdateProfile(ctx, userID, request)
	if err != nil {
		return ProfileResponse{}, err
	}

	return ProfileResponse{User: user}, nil
}

func (s Service) ChangePassword(ctx context.Context, userID uint64, request ChangePasswordRequest) error {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.CurrentPassword)); err != nil {
		return ErrInvalidCurrentPassword
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePasswordHash(ctx, userID, string(passwordHash))
}

func trimOptional(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
