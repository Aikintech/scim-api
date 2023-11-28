package validation

type StorePlaylistSchema struct {
	Title       string   `json:"title" validate:"required,min=10,max=30"`
	Description string   `json:"description" validate:"-"`
	Podcasts    []string `json:"podcasts" validate:"min=1"`
}
