package account

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"

	"github.com/CulipBlue/backend_ednic/internal/middleware"
	"github.com/CulipBlue/backend_ednic/internal/shared/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (h Handler) RegisterRoutes(router gin.IRoutes) {
	router.GET("/profile", h.getProfile)
	router.PATCH("/profile", h.updateProfile)
	router.PATCH("/password", h.changePassword)
}

// getProfile godoc
// @Summary Get account profile
// @Description Returns the authenticated user's account profile.
// @Tags Account
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ProfileSuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/account/profile [get]
func (h Handler) getProfile(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, 401, "Authentication required", nil)
		return
	}

	result, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, 404, "User not found", nil)
		return
	}

	response.OK(c, "Profile", result)
}

// updateProfile godoc
// @Summary Update account profile
// @Description Updates the authenticated user's profile.
// @Tags Account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Update profile request"
// @Success 200 {object} ProfileSuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/v1/account/profile [patch]
func (h Handler) updateProfile(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, 401, "Authentication required", nil)
		return
	}

	var request UpdateProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, 400, "Validation failed", gin.H{"request": err.Error()})
		return
	}

	result, err := h.service.UpdateProfile(c.Request.Context(), userID, request)
	if err != nil {
		if isDuplicateEntry(err) {
			response.Error(c, 409, "Username already used", nil)
			return
		}
		response.Error(c, 500, "Failed to update profile", gin.H{"error": err.Error()})
		return
	}

	response.OK(c, "Profile updated successfully", result)
}

// changePassword godoc
// @Summary Change account password
// @Description Changes the authenticated user's password after validating the current password.
// @Tags Account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Change password request"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/account/password [patch]
func (h Handler) changePassword(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, 401, "Authentication required", nil)
		return
	}

	var request ChangePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, 400, "Validation failed", gin.H{"request": err.Error()})
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID, request); err != nil {
		if errors.Is(err, ErrInvalidCurrentPassword) {
			response.Error(c, 401, "Current password is invalid", nil)
			return
		}
		response.Error(c, 500, "Failed to change password", gin.H{"error": err.Error()})
		return
	}

	response.OK(c, "Password changed successfully", nil)
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}

	return strings.Contains(strings.ToLower(err.Error()), "duplicate")
}
