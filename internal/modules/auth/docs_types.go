package auth

type AuthSuccessResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"Logged in successfully"`
	Data    AuthResponse `json:"data"`
	Errors  any          `json:"errors" swaggertype:"object"`
}

type UserSuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"User profile"`
	Data    User   `json:"data"`
	Errors  any    `json:"errors" swaggertype:"object"`
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Validation failed"`
	Data    any    `json:"data" swaggertype:"object"`
	Errors  any    `json:"errors" swaggertype:"object"`
}
