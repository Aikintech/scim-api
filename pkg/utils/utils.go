package utils

import (
	"crypto/rand"
	"io"
	"os"
	"time"

	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateCoe(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func GenerateUserToken(user *models.User, tokenType string, reference string) (string, error) {
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
	result := config.DB.Model(&models.UserToken{}).Create(&models.UserToken{
		UserID:      user.ID,
		Reference:   reference,
		Token:       t,
		Whitelisted: true,
	})
	if result.Error != nil {
		return "", result.Error
	}

	return t, nil
}
