package definitions

type ValidationErr struct {
	Field   string   `json:"field"`
	Reasons []string `json:"reasons"`
}

type MessageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ValidationErrsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Errors  []ValidationErr
}

type DataResponse[T any] struct {
	Code int `json:"code"`
	Data T   `json:"data"`
}

type Token struct {
	Reference string
	Token     string
}
