package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/CulipBlue/backend_ednic/internal/config"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInactiveUser       = errors.New("user is not active")
)

type Claims struct {
	UserID uint64 `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type Service struct {
	cfg  config.Config
	repo Repository
}

func NewService(cfg config.Config, repo Repository) Service {
	return Service{cfg: cfg, repo: repo}
}

func (s Service) Register(ctx context.Context, request RegisterRequest) (AuthResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, err
	}

	user, err := s.repo.CreateUser(ctx, User{
		Name:         strings.TrimSpace(request.Name),
		Username:     strings.TrimSpace(request.Username),
		Email:        strings.ToLower(strings.TrimSpace(request.Email)),
		PasswordHash: string(passwordHash),
		Role:         RoleUser,
		Status:       StatusActive,
	})
	if err != nil {
		return AuthResponse{}, err
	}

	token, err := s.generateToken(user)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{Token: token, User: user}, nil
}

func (s Service) CreateStaffUser(ctx context.Context, request CreateStaffUserRequest) (User, error) {
	if request.Role != RoleAdmin && request.Role != RoleSuperAdmin {
		request.Role = RoleAdmin
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	return s.repo.CreateUser(ctx, User{
		Name:         strings.TrimSpace(request.Name),
		Username:     strings.TrimSpace(request.Username),
		Email:        strings.ToLower(strings.TrimSpace(request.Email)),
		PasswordHash: string(passwordHash),
		Role:         request.Role,
		Status:       StatusActive,
	})
}

func (s Service) Login(ctx context.Context, request LoginRequest) (AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(request.Email)))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return AuthResponse{}, ErrInvalidCredentials
		}
		return AuthResponse{}, err
	}

	if user.Status != StatusActive {
		return AuthResponse{}, ErrInactiveUser
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		return AuthResponse{}, ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{Token: token, User: user}, nil
}

func (s Service) FindUserByID(ctx context.Context, userID uint64) (User, error) {
	return s.repo.FindByID(ctx, userID)
}

func (s Service) ParseToken(tokenValue string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func (s Service) generateToken(user User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: user.ID,
		Role:   user.Role,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(s.cfg.JWTAccessTokenTTLMinutes) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
