package definitions

type ResetPasswordRequest struct {
	Code                 string `json:"code"`
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
}

type StorePrayerRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required,len=10"`
	CountryCode string `json:"countryCode" validate:"required,len=2"`
}

type StorePostRequest struct {
	Title           string  `json:"title" validate:"required,min=3"`
	Body            string  `json:"body" validate:"required,min=3"`
	Published       bool    `json:"published" validate:"required"`
	ExcerptImageURL string  `json:"excerptImage" validate:"omitnil,validfilekey"`
	MinutesToRead   float32 `json:"minutesToRead" validate:"required"`
}

type StoreCommentRequest struct {
	Comment string `json:"comment" validate:"required,min=3,max=400"`
}
