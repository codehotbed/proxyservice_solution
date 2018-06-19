package utils

import (
	"ProxyService/proxyservice/openapi/autogen/proxyservice/models"
)

// NewErrorResponse Creates ErrorResponse json model and fills it with given parameters
func NewErrorResponse(code int64, message string) *models.ErrorResponse {
	return &models.ErrorResponse{Code: code, Message: message}
}
