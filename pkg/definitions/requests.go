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
	StartDateTime   string `json:"startDateTime" validate:"required,dateformat=2020-10-10 10:10:10"`
	EndDateTime     string `json:"endDateTime" validate:"required,dateformat=2020-10-10 10:10:10"`
	ExcerptImageURL string `json:"excerptImage" validate:"omitnil,validfilekey"`
}

type PlaylistPodcastsRequest struct {
	Podcasts []string `json:"podcasts" validate:"required,min=1"`
}
