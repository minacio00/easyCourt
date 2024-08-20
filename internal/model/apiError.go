package model

// APIError example used for Swagger response documentation.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
