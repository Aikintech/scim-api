package definitions

type MessageResponse struct {
	Message string `json:"message"`
}

type ValidationErrsResponse struct {
	Message string `json:"message"`
	Errors  []ValidationErr
}

type DataResponse[T any] struct {
	Code    int     `json:"code"`
	Data    T       `json:"data"`
	Message *string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
