package utils

import (
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateCoe(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(crand.Reader, b, max)
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

func GenerateRandomString(length int) string {
	pool := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := ""
	for i := 0; i < length; i++ {
		result += string(pool[rand.Intn(len(pool))])
	}
	return result
}

func DumpRoutesToFile(app *fiber.App) error {
	routes := app.GetRoutes(true)

	// Open the file for writing
	files := []string{"routes.txt", "routes.json"}
	error := error(nil)

	for _, filename := range files {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		if filename == "routes.txt" {
			for _, route := range routes {
				_, err := file.WriteString(fmt.Sprintf("%s %s\n", route.Method, route.Path))
				if err != nil {
					return err
				}
			}
		}

		if filename == "routes.json" {
			data, _ := json.MarshalIndent(routes, "", "  ")

			_, err = file.Write(data)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Routes dumped to %s\n", filename)
	}

	return error
}

func GetMimeExtension(mime string) string {
	split := strings.Split(mime, "/")
	if len(split) > 1 {
		return split[len(split)-1]
	}

	return ""
}
