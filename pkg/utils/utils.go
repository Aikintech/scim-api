package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/models"
	validationschemas "github.com/aikintech/scim/pkg/validation"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nedpals/supabase-go"

	"github.com/go-playground/validator/v10"
	"github.com/matthewhartstonge/argon2"
)

func ValidateStruct(schema interface{}) []definitions.ValidationErr {
	var errs []definitions.ValidationErr
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Custom error validation registration
	err := validate.RegisterValidation("isValidPassword", validationschemas.IsValidPasswordValidation)

	if err != nil {
		fmt.Println("Error registering custom validation :", err.Error())
	}

	// Validate struct
	err = validate.Struct(schema)

	if err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)

		for i, err := range validationErrors {
			field := strings.ToLower(err.Field())

			if exists := existsInValidationErrs(field, errs); exists != false {
				errs[i].Reasons = append(errs[i].Reasons, getValidationMessage(err))
			} else {
				errs = append(errs, definitions.ValidationErr{Field: field, Reasons: []string{getValidationMessage(err)}})
			}
		}
	}

	return errs
}

func MakePasswordHash(password string) (string, error) {
	argon := argon2.DefaultConfig()

	encoded, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return "", errors.New("Error hashing password")
	}

	return string(encoded), nil
}

func VerifyPasswordHash(password string, hashed string) (bool, error) {
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(hashed))
	if err != nil {
		return false, errors.New("Error verifying password")
	}

	return ok, nil
}

func GenerateUserToken(user *models.User, tokenType string) (string, error) {
	// Create the Claims
	expiry := time.Now().Add(time.Hour * 1).Unix()
	if tokenType == "refresh" {
		expiry = time.Now().Add(time.Hour * 24).Unix()
	}
	claims := jwt.MapClaims{
		"sub":       user.ID,
		"tokenType": tokenType,
		"exp":       expiry,
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("APP_KEY")))

	return t, err
}

/** Helpers **/

// getValidationMessage generates a human-readable validation message for a given validation error
func getValidationMessage(err validator.FieldError) string {
	fieldName := strings.ToLower(err.Field())

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required.", fieldName)
	case "min":
		return fmt.Sprintf("The %s field must be at least %s.", fieldName, err.Param())
	case "max":
		return fmt.Sprintf("The %s field must be at most %s.", fieldName, err.Param())
	case "len":
		return fmt.Sprintf("The %s field must have a length of %s.", fieldName, err.Param())
	case "email":
		return fmt.Sprintf("The %s field must be a valid email.", fieldName)
	case "url":
		return fmt.Sprintf("The %s field must be a valid URL.", fieldName)
	case "numeric":
		return fmt.Sprintf("The %s field must be a numeric value.", fieldName)
	case "alpha":
		return fmt.Sprintf("The %s field must contain only alphabetic characters.", fieldName)
	case "oneof":
		{
			split := strings.Split(err.Param(), " ")
			joined := strings.Join(split, ", ")

			return fmt.Sprintf("The %s field must be one of the following: %s.", fieldName, joined)
		}
	case "isValidPassword":
		return fmt.Sprintf("The %s field must contain at least one uppercase, one lowercase, one number and one special case character.", fieldName)

	// Add more cases for other validation tags as needed
	default:
		return fmt.Sprintf("Validation failed for %s with tag %s.", fieldName, err.Tag())
	}
}

// existsInValidationErrs checks if a field exists in a slice of ValidationErr
func existsInValidationErrs(field string, errs []definitions.ValidationErr) bool {
	result := false

	for _, key := range errs {
		if key.Field == field {
			result = true
			break
		}
	}

	return result
}

func LoginSupabaseUser() {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	supabaseClient := supabase.CreateClient(supabaseURL, supabaseKey)

	ctx := context.Background()
	user, err := supabaseClient.Auth.SignIn(ctx, supabase.UserCredentials{
		Email:    "nanaaikinson24@gmail.com",
		Password: "password",
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(user.User.ID)

	fmt.Println(user.AccessToken)
}
