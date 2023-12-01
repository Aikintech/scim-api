package definitions

type ValidationErr struct {
	Field   string   `json:"field"`
	Reasons []string `json:"reasons"`
}

type Token struct {
	Reference string
	Token     string
}

/**** Requests ****/
type ResetPasswordRequest struct {
	Code                 string `json:"code"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=Password"`
}

/**** Responses ****/
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
