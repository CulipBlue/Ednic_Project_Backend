package main

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/CulipBlue/backend_ednic/internal/config"
	"github.com/CulipBlue/backend_ednic/internal/database"
	"github.com/CulipBlue/backend_ednic/internal/modules/auth"
)

func main() {
	cfg := config.Load()
	if err := validateConfig(cfg); err != nil {
		log.Fatalf("invalid bootstrap config: %v", err)
	}

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	repo := auth.NewRepository(db)
	service := auth.NewService(cfg, repo)

	existing, err := repo.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(cfg.SuperAdminEmail)))
	if err == nil {
		log.Printf("super admin email already exists with role=%s username=%s", existing.Role, existing.Username)
		return
	}
	if !errors.Is(err, auth.ErrUserNotFound) {
		log.Fatalf("failed to check existing super admin: %v", err)
	}

	user, err := service.CreateStaffUser(ctx, auth.CreateStaffUserRequest{
		Name:     cfg.SuperAdminName,
		Username: cfg.SuperAdminUsername,
		Email:    cfg.SuperAdminEmail,
		Password: cfg.SuperAdminPassword,
		Role:     auth.RoleSuperAdmin,
	})
	if err != nil {
		log.Fatalf("failed to create super admin: %v", err)
	}

	log.Printf("super admin created successfully: id=%d email=%s username=%s", user.ID, user.Email, user.Username)
}

func validateConfig(cfg config.Config) error {
	missing := make([]string, 0)
	if strings.TrimSpace(cfg.SuperAdminName) == "" {
		missing = append(missing, "SUPER_ADMIN_NAME")
	}
	if strings.TrimSpace(cfg.SuperAdminUsername) == "" {
		missing = append(missing, "SUPER_ADMIN_USERNAME")
	}
	if strings.TrimSpace(cfg.SuperAdminEmail) == "" {
		missing = append(missing, "SUPER_ADMIN_EMAIL")
	}
	if strings.TrimSpace(cfg.SuperAdminPassword) == "" {
		missing = append(missing, "SUPER_ADMIN_PASSWORD")
	}
	if len(missing) > 0 {
		return errors.New("missing " + strings.Join(missing, ", "))
	}

	return nil
}
