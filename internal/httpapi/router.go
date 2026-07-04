package httpapi

import (
	"context"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/CulipBlue/backend_ednic/internal/config"
	"github.com/CulipBlue/backend_ednic/internal/middleware"
	"github.com/CulipBlue/backend_ednic/internal/modules/account"
	"github.com/CulipBlue/backend_ednic/internal/modules/auth"
	"github.com/CulipBlue/backend_ednic/internal/modules/products"
	"github.com/CulipBlue/backend_ednic/internal/shared/response"
)

type Server struct {
	cfg config.Config
	db  *sql.DB
}

func NewRouter(cfg config.Config, db *sql.DB) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	server := Server{cfg: cfg, db: db}
	router := gin.New()
	_ = router.SetTrustedProxies(nil)
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", server.health)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(cfg, authRepo)
	authHandler := auth.NewHandler(authService)
	accountRepo := account.NewRepository(db)
	accountService := account.NewService(accountRepo)
	accountHandler := account.NewHandler(accountService)
	productRepo := products.NewRepository(db)
	productService := products.NewService(productRepo)
	productHandler := products.NewHandler(productService)

	v1 := router.Group("/api/v1")
	v1.GET("/health/db", server.databaseHealth)
	productHandler.RegisterPublicRoutes(v1)

	authRoutes := v1.Group("/auth")
	authHandler.RegisterRoutes(authRoutes)

	authProtectedRoutes := v1.Group("/auth")
	authProtectedRoutes.Use(middleware.RequireAuth(authService))
	authProtectedRoutes.GET("/me", server.me(authService))

	accountRoutes := v1.Group("/account")
	accountRoutes.Use(middleware.RequireAuth(authService))
	accountHandler.RegisterRoutes(accountRoutes)

	adminRoutes := v1.Group("/admin")
	adminRoutes.Use(middleware.RequireAuth(authService), middleware.RequireAdmin())
	adminRoutes.GET("/health", server.adminHealth)
	productHandler.RegisterAdminRoutes(adminRoutes)

	return router
}

// health godoc
// @Summary API health check
// @Description Verifies the API process is running.
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (s Server) health(c *gin.Context) {
	response.OK(c, "API is running", gin.H{
		"service": "ednic-backend",
		"env":     s.cfg.AppEnv,
	})
}

// databaseHealth godoc
// @Summary Database health check
// @Description Verifies the API can connect to MySQL.
// @Tags Health
// @Produce json
// @Success 200 {object} DatabaseHealthResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/v1/health/db [get]
func (s Server) databaseHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		response.Error(c, 503, "Database connection failed", gin.H{
			"error": err.Error(),
		})
		return
	}

	response.OK(c, "Database connection is healthy", gin.H{
		"database": "mysql",
	})
}

// me godoc
// @Summary Get current user profile
// @Description Returns the authenticated user profile.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} auth.UserSuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/auth/me [get]
func (s Server) me(authService auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.CurrentUserID(c)
		if !ok {
			response.Error(c, 401, "Authentication required", nil)
			return
		}

		user, err := authService.FindUserByID(c.Request.Context(), userID)
		if err != nil {
			response.Error(c, 404, "User not found", nil)
			return
		}

		response.OK(c, "User profile", user)
	}
}

// adminHealth godoc
// @Summary Admin health check
// @Description Verifies the authenticated user has admin access.
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} AdminHealthResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /api/v1/admin/health [get]
func (s Server) adminHealth(c *gin.Context) {
	response.OK(c, "Admin access granted", gin.H{
		"scope": "admin",
	})
}
