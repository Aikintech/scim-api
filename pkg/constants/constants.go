package constants

const (
	PODCAST_URL        = "https://www.podcasts.com/rss_feed/00b03551b395628591a24c0ab6050926"
	JWT_CONTEXT_KEY    = "jwt-claims"
	USER_CONTEXT_KEY   = "user"
	PODCASTS_CACHE_KEY = "podcasts"
)

var (
	UPLOAD_TYPES       = []string{"testimony", "excerpt", "avatar"}
	ALLOWED_MIME_TYPES = []string{"png", "jpeg", "jpg", "mov", "mp4"}
)
