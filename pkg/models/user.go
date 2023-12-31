package models

import (
	"fmt"
	"os"
	"time"

	"github.com/aikintech/scim-api/pkg/utils"
	mapSet "github.com/deckarep/golang-set/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID              string `gorm:"primaryKey;size:40"`
	ExternalID      string `gorm:"size:40"`
	FirstName       string `gorm:"not null"`
	LastName        string `gorm:"not null"`
	Email           string `gorm:"not null;index"`
	Password        string
	EmailVerifiedAt *time.Time
	SignUpProvider  string `gorm:"not null"`
	Avatar          string
	PhoneNumber     string `gorm:"size:40"`
	Channels        datatypes.JSON
	CreatedAt       time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime;not null"`

	// Relations
	Playlists         []*Playlist
	PrayerRequests    []*PrayerRequest
	UserTokens        []*UserToken
	VerificationCodes []*VerificationCode
	Posts             []*Post
	Comments          []*Comment
	Roles             []*Role       `gorm:"many2many:role_user"`
	Permissions       []*Permission `gorm:"many2many:permission_user"`
}

type AuthUserResource struct {
	ID            string                `json:"id"`
	FirstName     string                `json:"firstName"`
	LastName      string                `json:"lastName"`
	Email         string                `json:"email"`
	EmailVerified bool                  `json:"emailVerified"`
	Avatar        string                `json:"avatar"`
	AvatarKey     string                `json:"avatarKey"`
	Channels      datatypes.JSON        `json:"channels"`
	Permissions   []*PermissionResource `json:"permissions"`
}

type BackofficeUser struct {
	ID            string         `json:"id"`
	FirstName     string         `json:"firstName"`
	LastName      string         `json:"lastName"`
	Email         string         `json:"email"`
	EmailVerified bool           `json:"emailVerified"`
	Avatar        string         `json:"avatar"`
	Channels      datatypes.JSON `json:"channels"`
	Role          string         `json:"role"`
}

type BackofficeUserFull struct {
	ID            string         `json:"id"`
	FirstName     string         `json:"firstName"`
	LastName      string         `json:"lastName"`
	Email         string         `json:"email"`
	EmailVerified bool           `json:"emailVerified"`
	Avatar        string         `json:"avatar"`
	Channels      datatypes.JSON `json:"channels"`
	Role          *RoleResource  `json:"role"`
}

type UserRel struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
}

func (u *User) BeforeCreate(*gorm.DB) error {
	u.ID = ulid.Make().String()

	return nil
}

func ToAuthUserResource(u *User) AuthUserResource {
	// Generate avatarURL
	avatar, err := utils.GenerateS3FileURL(u.Avatar)
	if err != nil {
		fmt.Println("Error generating avatar url", err.Error())
	}

	// Get user permissions
	permissions := mapSet.NewSet[*PermissionResource]()
	if len(u.Roles) > 0 {
		for _, p := range u.Roles[0].Permissions {
			permissions.Add(PermissionToResource(p))
		}
	}
	if len(u.Permissions) > 0 {
		for _, p := range u.Permissions {
			permissions.Add(PermissionToResource(p))
		}
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
		Permissions:   permissions.ToSlice(),
	}
}

func ToUserRelResource(u *User) UserRel {
	// Generate avatarURL
	avatar, err := utils.GenerateS3FileURL(u.Avatar)
	if err != nil {
		fmt.Println("Error generating avatar url", err.Error())
	}

	return UserRel{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Avatar:    avatar,
	}
}

func ToBackofficeUserResource(u *User) *BackofficeUser {
	var role *Role

	// Generate avatarURL
	avatar, err := utils.GenerateS3FileURL(u.Avatar)
	if err != nil {
		fmt.Println("Error generating avatar url", err.Error())
	}

	if len(u.Roles) > 0 {
		role = u.Roles[0]
	}

	return &BackofficeUser{
		ID:            u.ID,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Email:         u.Email,
		EmailVerified: u.EmailVerifiedAt != nil,
		Avatar:        avatar,
		Channels:      u.Channels,
		Role:          role.DisplayName,
	}
}

func ToBackofficeUserFullResource(u *User) *BackofficeUserFull {
	// Generate avatarURL
	avatar, err := utils.GenerateS3FileURL(u.Avatar)
	if err != nil {
		fmt.Println("Error generating avatar url", err.Error())
	}

	role := u.Roles[0]

	return &BackofficeUserFull{
		ID:            u.ID,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Email:         u.Email,
		EmailVerified: u.EmailVerifiedAt != nil,
		Avatar:        avatar,
		Channels:      u.Channels,
		Role:          RoleToResource(role),
	}
}

func UsersToBackofficeUserResourceCollection(users []*User) []*BackofficeUser {
	var resources []*BackofficeUser

	for _, u := range users {
		resources = append(resources, ToBackofficeUserResource(u))
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
