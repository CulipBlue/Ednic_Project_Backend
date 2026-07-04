package auth

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"

	"github.com/CulipBlue/backend_ednic/internal/shared/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (h Handler) RegisterRoutes(router gin.IRoutes) {
	router.POST("/register", h.register)
	router.POST("/login", h.login)
}

// register godoc
// @Summary Register user
// @Description Creates a new user account with role user.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request"
// @Success 200 {object} AuthSuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/register [post]
func (h Handler) register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, 400, "Validation failed", gin.H{"request": err.Error()})
		return
	}

	result, err := h.service.Register(c.Request.Context(), request)
	if err != nil {
		if isDuplicateEntry(err) {
			response.Error(c, 409, "Email or username already registered", nil)
			return
		}
		response.Error(c, 500, "Failed to register user", gin.H{"error": err.Error()})
		return
	}

	response.OK(c, "Registered successfully", result)
}

// login godoc
// @Summary Login user
// @Description Authenticates a user and returns a JWT access token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} AuthSuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func (h Handler) login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, 400, "Validation failed", gin.H{"request": err.Error()})
		return
	}

	result, err := h.service.Login(c.Request.Context(), request)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) || errors.Is(err, ErrInactiveUser) {
			response.Error(c, 401, "Invalid email or password", nil)
			return
		}
		response.Error(c, 500, "Failed to login", gin.H{"error": err.Error()})
		return
	}

	response.OK(c, "Logged in successfully", result)
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}

	return strings.Contains(strings.ToLower(err.Error()), "duplicate")
}
