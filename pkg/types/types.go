package types

type ValidationError struct {
	Field  string   `json:"field"`
	Errors []string `json:"errors"`
}
