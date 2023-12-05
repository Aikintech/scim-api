package validation

type LoginSchema struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Channel  string `json:"channel" validate:"required,oneof=web mobile"`
}

type RegisterSchema struct {
	FirstName string `json:"firstName" validate:"required,min=3,max=40"`
	LastName  string `json:"lastName" validate:"required,min=3,max=40"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=40,validpassword"`
	Channel   string `json:"channel" validate:"required,oneof=web mobile"`
}

type SocialAuthSchema struct {
	UserToken string `json:"userToken" validate:"required"`
}

type EmailVerificationSchema struct {
	Email string `json:"email" validate:"required,email"`
}
