package constants

import "os"

const (
	PODCAST_URL        = "https://www.podcasts.com/rss_feed/00b03551b395628591a24c0ab6050926"
	JWT_CONTEXT_KEY    = "jwt-claims"
	USER_CONTEXT_KEY   = "user"
	PODCASTS_CACHE_KEY = "podcasts"
	DATE_TIME_FORMAT   = "2023-12-21T21:00:00"
)

var (
	UPLOAD_TYPES       = []string{"testimony", "excerpt", "avatar"}
	ALLOWED_MIME_TYPES = []string{"png", "jpeg", "jpg", "mov", "mp4"}

	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_DATABASE = os.Getenv("DB_DATABASE")
	DB_USERNAME = os.Getenv("DB_USERNAME")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
)
