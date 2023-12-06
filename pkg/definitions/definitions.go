package definitions

type ValidationErr struct {
	Field   string   `json:"field"`
	Reasons []string `json:"reasons"`
}

type Token struct {
	Reference string
	Token     string
}

type Map map[string]interface{}

/**** Requests ****/
type ResetPasswordRequest struct {
	Code                 string `json:"code"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,validpassword"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=Password"`
}

type StorePlaylistRequest struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description"`
	Podcasts    []string `json:"podcasts" validate:"required,min=1"`
}

type StoreEventRequest struct {
	Title           string `json:"title" validate:"required,min:20"`
	Description     string `json:"description" validate:"required,min:20"`
	Location        string `json:"location" validate:"required"`
	StartDateTime   string `json:"startDateTime" validate:"required"`
	EndDateTime     string `json:"endDateTime" validate:"required"`
	ExcerptImageURL string `json:"excerptImage" validate:"-"`
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
	Code    int     `json:"code"`
	Data    T       `json:"data"`
	Message *string `json:"message"`
}
