package products

type ProductSuccessResponse struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message" example:"Product"`
	Data    ProductResponse `json:"data"`
	Errors  any             `json:"errors" swaggertype:"object"`
}

type ProductListSuccessResponse struct {
	Success bool                 `json:"success" example:"true"`
	Message string               `json:"message" example:"Products"`
	Data    ListProductsResponse `json:"data"`
	Errors  any                  `json:"errors" swaggertype:"object"`
}

type CategorySuccessResponse struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"Category"`
	Data    CategoryResponse `json:"data"`
	Errors  any              `json:"errors" swaggertype:"object"`
}

type CategoryListSuccessResponse struct {
	Success bool       `json:"success" example:"true"`
	Message string     `json:"message" example:"Categories"`
	Data    []Category `json:"data"`
	Errors  any        `json:"errors" swaggertype:"object"`
}

type MessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Product deleted successfully"`
	Data    any    `json:"data" swaggertype:"object"`
	Errors  any    `json:"errors" swaggertype:"object"`
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Validation failed"`
	Data    any    `json:"data" swaggertype:"object"`
	Errors  any    `json:"errors" swaggertype:"object"`
}
