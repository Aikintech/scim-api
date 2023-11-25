package types

type FormattedValidationErrs []ValidationErr
type GookitErrs map[string]map[string]string

type ValidationErr struct {
	Field   string   `json:"field"`
	Reasons []string `json:"reasons"`
}

// Response types
type MessageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type ValidationErrorResponse struct {
	Errors []ValidationErr `json:"errors"`
}

type DataResponse[T any] struct {
	MessageResponse
	Data T `json:"data"`
}
