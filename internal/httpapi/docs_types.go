package httpapi

type HealthData struct {
	Service string `json:"service" example:"ednic-backend"`
	Env     string `json:"env" example:"local"`
}

type DatabaseHealthData struct {
	Database string `json:"database" example:"mysql"`
}

type AdminHealthData struct {
	Scope string `json:"scope" example:"admin"`
}

type HealthResponse struct {
	Success bool       `json:"success" example:"true"`
	Message string     `json:"message" example:"API is running"`
	Data    HealthData `json:"data"`
	Errors  any        `json:"errors" swaggertype:"object"`
}

type DatabaseHealthResponse struct {
	Success bool               `json:"success" example:"true"`
	Message string             `json:"message" example:"Database connection is healthy"`
	Data    DatabaseHealthData `json:"data"`
	Errors  any                `json:"errors" swaggertype:"object"`
}

type AdminHealthResponse struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message" example:"Admin access granted"`
	Data    AdminHealthData `json:"data"`
	Errors  any             `json:"errors" swaggertype:"object"`
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Validation failed"`
	Data    any    `json:"data" swaggertype:"object"`
	Errors  any    `json:"errors" swaggertype:"object"`
}
