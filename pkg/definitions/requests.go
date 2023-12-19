package definitions

type ResetPasswordRequest struct {
	Key                  string `json:"key" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,validpassword"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=Password"`
}

type StorePlaylistRequest struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description"`
	Podcasts    []string `json:"podcasts"`
}

type StoreEventRequest struct {
	Title           string `json:"title" validate:"required,min=3"`
	Description     string `json:"description" validate:"required,min=3"`
	Location        string `json:"location" validate:"required"`
	Published       bool   `json:"published" validate:"boolean"`
	StartDateTime   string `json:"startDateTime" validate:"required,dateformat=2006-01-02T15:04:05.000Z"`
	EndDateTime     string `json:"endDateTime" validate:"required,dateformat=2006-01-02T15:04:05.000Z"`
	ExcerptImageURL string `json:"excerptImage" validate:"omitnil,validfilekey"`
}

type PlaylistPodcastsRequest struct {
	Podcasts []string `json:"podcasts" validate:"required,min=1"`
}

type SocialAuthRequest struct {
	Provider string `json:"provider" validate:"required,oneof=apple google"`
	Token    string `json:"token" validate:"required"`
	Channel  string `json:"channel" validate:"required,oneof=mobile web"`
}

type StorePrayerRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required,len=10"`
	CountryCode string `json:"countryCode" validate:"required,len=2"`
}

type StorePostRequest struct {
	Title          string  `json:"title" validate:"required,min=3"`
	Body           string  `json:"body" validate:"required,min=3"`
	Published      bool    `json:"published" validate:"boolean"`
	ExcerptImage   string  `json:"excerptImage" validate:"omitnil,validfilekey"`
	MinutesToRead  float32 `json:"minutesToRead" validate:"required"`
	IsAnnouncement bool    `json:"isAnnouncement" validate:"boolean"`
}

type StoreCommentRequest struct {
	Comment string `json:"comment" validate:"required,min=3,max=400"`
}

type VerifyEmailRequest struct {
	Code   string `json:"code" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	Action string `json:"action" validate:"required,oneof=reset_password account_verification"`
}

type UpdateAvatarRequest struct {
	AvatarKey string `json:"avatarKey" validate:"required,validfilekey"`
	Action    string `json:"action" validate:"required,oneof=update remove"`
}

type UpdateUserDetailsRequest struct {
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	PhoneNumber string `json:"phoneNumber"`
	CountryCode string `json:"countryCode"`
}

type VerifyAccountRequest struct {
	Key   string `json:"key" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type TransactRequest struct {
	Amount         float32 `json:"amount" validate:"required,min=1"`
	IdempotencyKey string  `json:"idempotencyKey" validate:"required,min=26,max=36"`
	Type           string  `json:"type" validate:"required,oneof=tithe pledge offertory freewill other busing covenant_partner"`
	Currency       string  `json:"currency" validate:"required,oneof=USD GHS EUR GBP"`
	AccountNumber  string  `json:"accountNumber" validate:"required"`
	Method         string  `json:"method" validate:"required,oneof=mobile_money card"`
}
