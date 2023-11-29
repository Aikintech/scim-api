package validation

type CommentPodcastSchema struct {
	Comment string `json:"comment" validate:"required,min=3,max=400"`
}

type StorePlaylistSchema struct {
	Title       string   `json:"title" validate:"required,min=10,max=30"`
	Description string   `json:"description" validate:"-"`
	Podcasts    []string `json:"podcasts" validate:"min=1"`
}
