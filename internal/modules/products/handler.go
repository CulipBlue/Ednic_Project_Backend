package products

import (
	"errors"
	"net/http"
	"strconv"

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

func (h Handler) RegisterPublicRoutes(router gin.IRoutes) {
	router.GET("/products", h.listPublicProducts)
	router.GET("/products/:slug", h.getPublicProduct)
	router.GET("/product-categories", h.listPublicCategories)
}

func (h Handler) RegisterAdminRoutes(router gin.IRoutes) {
	router.GET("/products", h.listAdminProducts)
	router.POST("/products", h.createProduct)
	router.GET("/products/:id", h.getAdminProduct)
	router.PUT("/products/:id", h.updateProduct)
	router.DELETE("/products/:id", h.deleteProduct)

	router.GET("/product-categories", h.listAdminCategories)
	router.POST("/product-categories", h.createCategory)
	router.PUT("/product-categories/:id", h.updateCategory)
	router.DELETE("/product-categories/:id", h.deleteCategory)
}

// listPublicProducts godoc
// @Summary List public products
// @Description Returns active products for public catalog.
// @Tags Products
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param search query string false "Search keyword"
// @Param category query string false "Category slug"
// @Success 200 {object} ProductListSuccessResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/products [get]
func (h Handler) listPublicProducts(c *gin.Context) {
	result, err := h.service.ListPublicProducts(c.Request.Context(), listQueryFromRequest(c, true))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list products", gin.H{"error": err.Error()})
		return
	}

	response.OK(c, "Products", result)
}

// getPublicProduct godoc
// @Summary Get public product detail
// @Description Returns active product detail by slug.
// @Tags Products
// @Produce json
// @Param slug path string true "Product slug"
// @Success 200 {object} ProductSuccessResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/products/{slug} [get]
func (h Handler) getPublicProduct(c *gin.Context) {
	result, err := h.service.FindPublicProductBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			response.Error(c, http.StatusNotFound, "Product not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get product", gin.H{"error": err.Error()})
		return
	}

	response.OK(c, "Product", result)
}

// listPublicCategories godoc
// @Summary List public product categories
// @Description Returns active product categories.
// @Tags Products
// @Produce json
// @Success 200 {object} CategoryListSuccessResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/product-categories [get]
func (h Handler) listPublicCategories(c *gin.Context) {
	result, err := h.service.ListPublicCategories(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list categories", gin.H{"error": err.Error()})
		return
	}

	response.OK(c, "Categories", result)
}

// listAdminProducts godoc
// @Summary List products for admin
// @Description Returns products for admin including draft and inactive items.
// @Tags Admin Products
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param search query string false "Search keyword"
// @Param status query string false "Status"
// @Param category query string false "Category slug"
// @Success 200 {object} ProductListSuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /api/v1/admin/products [get]
func (h Handler) listAdminProducts(c *gin.Context) {
	result, err := h.service.ListAdminProducts(c.Request.Context(), listQueryFromRequest(c, false))
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to list products"
		if errors.Is(err, ErrInvalidStatus) {
			status = http.StatusBadRequest
			message = "Invalid status"
		}
		response.Error(c, status, message, gin.H{"error": err.Error()})
		return
	}

	response.OK(c, "Products", result)
}

// createProduct godoc
// @Summary Create product
// @Description Creates a product from admin panel.
// @Tags Admin Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProductRequest true "Create product request"
// @Success 200 {object} ProductSuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/v1/admin/products [post]
func (h Handler) createProduct(c *gin.Context) {
	actorID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	var request CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation failed", gin.H{"request": err.Error()})
		return
	}

	result, err := h.service.CreateProduct(c.Request.Context(), request, actorID)
	if err != nil {
		h.handleWriteError(c, "Failed to create product", err)
		return
	}

	response.OK(c, "Product created successfully", result)
}

// getAdminProduct godoc
// @Summary Get product for admin
// @Description Returns product detail by id for admin.
// @Tags Admin Products
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} ProductSuccessResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/admin/products/{id} [get]
func (h Handler) getAdminProduct(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	result, err := h.service.FindAdminProductByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			response.Error(c, http.StatusNotFound, "Product not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to get product", gin.H{"error": err.Error()})
		return
	}

	response.OK(c, "Product", result)
}

// updateProduct godoc
// @Summary Update product
// @Description Updates product from admin panel.
// @Tags Admin Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body UpdateProductRequest true "Update product request"
// @Success 200 {object} ProductSuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/v1/admin/products/{id} [put]
func (h Handler) updateProduct(c *gin.Context) {
	actorID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Authentication required", nil)
		return
	}
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var request UpdateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation failed", gin.H{"request": err.Error()})
		return
	}

	result, err := h.service.UpdateProduct(c.Request.Context(), id, request, actorID)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			response.Error(c, http.StatusNotFound, "Product not found", nil)
			return
		}
		h.handleWriteError(c, "Failed to update product", err)
		return
	}

	response.OK(c, "Product updated successfully", result)
}

// deleteProduct godoc
// @Summary Delete product
// @Description Deletes product from admin panel.
// @Tags Admin Products
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/admin/products/{id} [delete]
func (h Handler) deleteProduct(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrProductNotFound) {
			response.Error(c, http.StatusNotFound, "Product not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete product", gin.H{"error": err.Error()})
		return
	}

	response.OK(c, "Product deleted successfully", nil)
}

// listAdminCategories godoc
// @Summary List product categories for admin
// @Description Returns all product categories including inactive categories.
// @Tags Admin Products
// @Produce json
// @Security BearerAuth
// @Success 200 {object} CategoryListSuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /api/v1/admin/product-categories [get]
func (h Handler) listAdminCategories(c *gin.Context) {
	result, err := h.service.ListAdminCategories(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list categories", gin.H{"error": err.Error()})
		return
	}

	response.OK(c, "Categories", result)
}

// createCategory godoc
// @Summary Create product category
// @Description Creates a product category from admin panel.
// @Tags Admin Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCategoryRequest true "Create category request"
// @Success 200 {object} CategorySuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/v1/admin/product-categories [post]
func (h Handler) createCategory(c *gin.Context) {
	var request CreateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation failed", gin.H{"request": err.Error()})
		return
	}

	result, err := h.service.CreateCategory(c.Request.Context(), request)
	if err != nil {
		h.handleWriteError(c, "Failed to create category", err)
		return
	}

	response.OK(c, "Category created successfully", result)
}

// updateCategory godoc
// @Summary Update product category
// @Description Updates a product category from admin panel.
// @Tags Admin Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param request body UpdateCategoryRequest true "Update category request"
// @Success 200 {object} CategorySuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/v1/admin/product-categories/{id} [put]
func (h Handler) updateCategory(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var request UpdateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation failed", gin.H{"request": err.Error()})
		return
	}

	result, err := h.service.UpdateCategory(c.Request.Context(), id, request)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, "Category not found", nil)
			return
		}
		h.handleWriteError(c, "Failed to update category", err)
		return
	}

	response.OK(c, "Category updated successfully", result)
}

// deleteCategory godoc
// @Summary Delete product category
// @Description Deletes a product category from admin panel.
// @Tags Admin Products
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/admin/product-categories/{id} [delete]
func (h Handler) deleteCategory(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.service.DeleteCategory(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, "Category not found", nil)
			return
		}
		h.handleWriteError(c, "Failed to delete category", err)
		return
	}

	response.OK(c, "Category deleted successfully", nil)
}

func (h Handler) handleWriteError(c *gin.Context, message string, err error) {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062:
			response.Error(c, http.StatusConflict, DuplicateMessage(err), nil)
			return
		case 1452:
			response.Error(c, http.StatusBadRequest, "Related data is invalid", gin.H{"error": ForeignKeyMessage(err)})
			return
		}
	}

	response.Error(c, http.StatusInternalServerError, message, gin.H{"error": err.Error()})
}

func listQueryFromRequest(c *gin.Context, public bool) ListProductsQuery {
	return ListProductsQuery{
		Page:     parseQueryInt(c, "page", 1),
		Limit:    parseQueryInt(c, "limit", 20),
		Search:   c.Query("search"),
		Status:   c.Query("status"),
		Category: c.Query("category"),
		Public:   public,
	}
}

func parseQueryInt(c *gin.Context, key string, fallback int) int {
	value := c.Query(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseIDParam(c *gin.Context, key string) (uint64, bool) {
	value := c.Param(key)
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil || parsed == 0 {
		response.Error(c, http.StatusBadRequest, "Invalid "+key, nil)
		return 0, false
	}
	return parsed, true
}
