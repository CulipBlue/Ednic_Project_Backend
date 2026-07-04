package products

import "time"

const (
	ProductStatusDraft    = "draft"
	ProductStatusActive   = "active"
	ProductStatusInactive = "inactive"

	CategoryStatusActive   = "active"
	CategoryStatusInactive = "inactive"
)

type Category struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Product struct {
	ID               uint64    `json:"id"`
	CategoryID       *uint64   `json:"category_id"`
	Category         *Category `json:"category,omitempty"`
	SKU              *string   `json:"sku"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	ShortDescription *string   `json:"short_description"`
	Description      *string   `json:"description"`
	Price            float64   `json:"price"`
	Stock            uint64    `json:"stock"`
	Material         *string   `json:"material"`
	WeightGram       *uint64   `json:"weight_gram"`
	LengthMM         *uint64   `json:"length_mm"`
	WidthMM          *uint64   `json:"width_mm"`
	HeightMM         *uint64   `json:"height_mm"`
	MainImageURL     *string   `json:"main_image_url"`
	Status           string    `json:"status"`
	CreatedBy        *uint64   `json:"created_by"`
	UpdatedBy        *uint64   `json:"updated_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ListProductsQuery struct {
	Page     int
	Limit    int
	Search   string
	Status   string
	Public   bool
	Category string
}

type ListProductsResponse struct {
	Items []Product `json:"items"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
	Total int64     `json:"total"`
}

type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=120"`
	Slug        string  `json:"slug" binding:"required,min=2,max=140"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Status      string  `json:"status" binding:"omitempty,oneof=active inactive"`
}

type UpdateCategoryRequest = CreateCategoryRequest

type CreateProductRequest struct {
	CategoryID       *uint64 `json:"category_id"`
	SKU              *string `json:"sku" binding:"omitempty,max=80"`
	Name             string  `json:"name" binding:"required,min=2,max=180"`
	Slug             string  `json:"slug" binding:"required,min=2,max=200"`
	ShortDescription *string `json:"short_description" binding:"omitempty,max=300"`
	Description      *string `json:"description"`
	Price            float64 `json:"price" binding:"required,gte=0"`
	Stock            uint64  `json:"stock"`
	Material         *string `json:"material" binding:"omitempty,max=120"`
	WeightGram       *uint64 `json:"weight_gram"`
	LengthMM         *uint64 `json:"length_mm"`
	WidthMM          *uint64 `json:"width_mm"`
	HeightMM         *uint64 `json:"height_mm"`
	MainImageURL     *string `json:"main_image_url" binding:"omitempty,max=500"`
	Status           string  `json:"status" binding:"omitempty,oneof=draft active inactive"`
}

type UpdateProductRequest = CreateProductRequest

type ProductResponse struct {
	Product Product `json:"product"`
}

type CategoryResponse struct {
	Category Category `json:"category"`
}
