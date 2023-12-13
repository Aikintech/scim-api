package models

import (
	"fmt"
	"os"
	"time"

	"github.com/aikintech/scim-api/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID                string `gorm:"primaryKey;size:40"`
	FirstName         string `gorm:"not null"`
	LastName          string `gorm:"not null"`
	Email             string `gorm:"not null;index"`
	Password          string
	EmailVerifiedAt   *time.Time
	SignUpProvider    string `gorm:"not null"`
	Avatar            string
	PhoneNumber       string
	Channels          datatypes.JSON
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Playlists         []*Playlist
	PrayerRequests    []*PrayerRequest
	UserTokens        []*UserToken
	VerificationCodes []*VerificationCode
	Posts             []*Post
	Comments          []*Comment
}

type AuthUserResource struct {
	ID            string         `json:"id"`
	FirstName     string         `json:"firstName"`
	LastName      string         `json:"lastName"`
	Email         string         `json:"email"`
	EmailVerified bool           `json:"emailVerified"`
	Avatar        string         `json:"avatar"`
	AvatarKey     string         `json:"avatarKey"`
	Channels      datatypes.JSON `json:"channels"`
}

func (u *User) BeforeCreate(*gorm.DB) error {
	u.ID = ulid.Make().String()

	return nil
}

func UserToResource(u *User) AuthUserResource {
	// Generate avatarURL
	avatar, err := utils.GenerateS3FileURL(u.Avatar)
	if err != nil {
		fmt.Println("Error generating avatar url", err.Error())
	}

	return AuthUserResource{
		ID:            u.ID,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Email:         u.Email,
		EmailVerified: u.EmailVerifiedAt != nil,
		Avatar:        avatar,
		AvatarKey:     u.Avatar,
		Channels:      u.Channels,
	}
}

func UsersToResourceCollection(users []*User) []AuthUserResource {
	var resources []AuthUserResource

	for _, user := range users {
		resources = append(resources, UserToResource(user))
	}

	return resources
}

func GenerateUserToken(user User, tokenType string, reference string) (string, error) {
	// Create the Claims
	expiry := time.Now().Add(time.Hour * 1).Unix()
	if tokenType == "refresh" {
		expiry = time.Now().Add(time.Hour * 24).Unix()
	}
	claims := jwt.MapClaims{
		"sub":       user.ID,
		"tokenType": tokenType,
		"reference": reference,
		"exp":       expiry,
		"iat":       time.Now().Unix(),
		"iss":       os.Getenv("APP_ISS"),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("APP_KEY")))
	if err != nil {
		return "", err
	}

	// Create user token
	// result := database.DB.Model(&UserToken{}).Create(&UserToken{
	// 	UserID:      user.ID,
	// 	Reference:   reference,
	// 	Token:       t,
	// 	Whitelisted: true,
	// })
	// if result.Error != nil {
	// 	return "", result.Error
	// }

	return t, nil
}
