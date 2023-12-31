package constants

import "os"

const (
	PODCAST_URL                             = "https://www.podcasts.com/rss_feed/00b03551b395628591a24c0ab6050926"
	JWT_CONTEXT_KEY                         = "jwt-claims"
	USER_CONTEXT_KEY                        = "user"
	PODCASTS_CACHE_KEY                      = "podcasts"
	DATE_TIME_FORMAT                        = "2006-01-02T15:04:05.000Z"
	ALPHABETS                               = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	MAILJET_VERIFY_MAIL_TEMPLATE_ID         = 5423172
	MAILJET_WELCOME_MAIL_TEMPLATE_ID        = 5423196
	MAILJET_RESET_PASSWORD_MAIL_TEMPLATE_ID = 5423195
	// NO_REPLY_EMAIL                          = "admin@scimapp.org"
	NO_REPLY_EMAIL                   = "noreply@scimapp.org"
	SUPPORT_EMAIL                    = "support@scimapp.org"
	USER_VERIFICATION_CODE_CACHE_KEY = "user_verification_code_" // + user.ID
	USER_CODE_ACTION_CACHE_KEY       = "user_code_action_"       // + user.ID
	YOUTUBE_CHANNEL_ID               = ""
	REGULAR_USER_ROLE                = "regular-user"
	SUPER_ADMIN_USER_ROLE            = "super-admin"
)

var (
	UPLOAD_TYPES       = []string{"testimony", "excerpt", "avatar"}
	ALLOWED_MIME_TYPES = []string{"png", "jpeg", "jpg", "mov", "mp4"}

	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_DATABASE = os.Getenv("DB_DATABASE")
	DB_USERNAME = os.Getenv("DB_USERNAME")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")

	PAYSTACK_PUBLIC_KEY = os.Getenv("PAYSTACK_PUBLIC_KEY")
	PAYSTACK_SECRET_KEY = os.Getenv("PAYSTACK_SECRET_KEY")
)
