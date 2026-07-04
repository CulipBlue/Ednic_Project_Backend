package account

type ProfileSuccessResponse struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message" example:"Profile updated successfully"`
	Data    ProfileResponse `json:"data"`
	Errors  any             `json:"errors" swaggertype:"object"`
}

type MessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Password changed successfully"`
	Data    any    `json:"data" swaggertype:"object"`
	Errors  any    `json:"errors" swaggertype:"object"`
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Validation failed"`
	Data    any    `json:"data" swaggertype:"object"`
	Errors  any    `json:"errors" swaggertype:"object"`
}
