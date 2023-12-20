package controllers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aikintech/scim-api/pkg/config"
	"github.com/aikintech/scim-api/pkg/database"
	"github.com/aikintech/scim-api/pkg/facades"
	"github.com/aikintech/scim-api/pkg/jobs"
	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/ttacon/libphonenumber"

	"github.com/aikintech/scim-api/pkg/constants"

	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/utils"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/golang-module/carbon/v2"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	// Parse request
	request := validation.LoginSchema{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Errors: errs,
		})
	}

	// Fetch User
	user := models.User{}
	result := database.DB.Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Invalid credentials provided",
		})
	}

	// Check password
	ok, err := utils.VerifyPasswordHash(request.Password, user.Password)
	if err != nil || !ok {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Invalid credentials provided",
		})
	}

	// Check if user is verified
	if user.EmailVerifiedAt == nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Account not verified. Please verify your email address",
		})
	}

	// Generate token
	reference := ulid.Make().String()
	accessToken, refreshToken, err := generateTokens(user, reference)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": models.UserToResource(&user),
		"tokens": fiber.Map{
			"access":  accessToken,
			"refresh": refreshToken,
		},
	})
}

func (a *AuthController) Register(c *fiber.Ctx) error {
	// Parse request
	request := validation.RegisterSchema{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{

			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Errors: errs,
		})
	}

	// Check if email exists
	user := models.User{}
	trx := database.DB.Begin()
	result := trx.Where("email = ?", request.Email).First(&user)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{

			Message: result.Error.Error(),
		})
	}
	if len(user.ID) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{

			Message: "An account with this email already exists",
		})
	}

	// Create user
	password, err := utils.MakePasswordHash(request.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}
	user.Email = request.Email
	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.Password = password
	user.Channels = datatypes.JSON([]byte(`["web", "mobile"]`))
	user.SignUpProvider = "local"
	user.EmailVerifiedAt = nil
	result = trx.Model(&models.User{}).Create(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: result.Error.Error(),
		})
	}

	// Save user verification code to redis
	code := utils.GenerateRandomNumbers(6)
	if err := config.RedisStore.Set(constants.USER_VERIFICATION_CODE_CACHE_KEY+user.ID, []byte(code), time.Minute*10); err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Send email verification mail
	go jobs.NewMail().SendUserVerificationMail(user, code)

	// Commit transaction
	trx.Commit()

	return c.Status(fiber.StatusCreated).JSON(definitions.SuccessResponse{
		Message: "Account created successfully",
	})
}

func (a *AuthController) RefreshToken(c *fiber.Ctx) error {
	reference := ulid.Make().String()
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	accessToken, err := models.GenerateUserToken(*user, "access", reference)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{

			Message: err.Error(),
		})
	}
	//
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access": accessToken,
	})
}

func (a *AuthController) ResendEmailVerification(c *fiber.Ctx) error {
	// Parse request
	request := validation.EmailVerificationSchema{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{

			Message: err.Error(),
		})
	}

	// Check if user exists
	message := ""
	user := new(models.User)
	result := database.DB.Model(&models.User{}).Where("email = ?", request.Email).First(&user)

	// Gorm error
	if result.Error != nil {
		message = result.Error.Error()

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = "An account with the selected email does not exist"
		}
	}

	// Sign up provider is not local
	if user.SignUpProvider != "local" {
		message = "Sorry you cannot reset your password with this sign up provider"
	}

	if len(message) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{

			Message: message,
		})
	}

	// Save user verification code to redis
	code := utils.GenerateRandomNumbers(6)
	if err := config.RedisStore.Set(constants.USER_VERIFICATION_CODE_CACHE_KEY+user.ID, []byte(code), time.Minute*10); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Send email verification mail
	go jobs.NewMail().SendUserVerificationMail(*user, code)

	return c.JSON(definitions.SuccessResponse{
		Message: "Email verification sent",
	})
}

func (a *AuthController) ForgotPassword(c *fiber.Ctx) error {
	// Parse request
	request := validation.EmailVerificationSchema{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Errors: errs,
		})
	}

	// Find user
	user := models.User{}
	result := database.DB.Model(&models.User{}).Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		message := "No account is associated with the email provided"

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = result.Error.Error()
		}

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{

			Message: message,
		})
	}

	// Check user's sign up provider
	if user.SignUpProvider != "local" {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{

			Message: "The sign up provider for this account does not support password reset",
		})
	}

	// Save user verification code to redis
	code := utils.GenerateRandomNumbers(6)
	if err := config.RedisStore.Set(constants.USER_VERIFICATION_CODE_CACHE_KEY+user.ID, []byte(code), time.Minute*10); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Send email verification mail
	go jobs.NewMail().SendUserPasswordResetMail(user, code)

	return c.JSON(definitions.MessageResponse{
		Message: "Password reset email sent",
	})
}

func (a *AuthController) ResetPassword(c *fiber.Ctx) error {
	// Parse request
	request := definitions.ResetPasswordRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{

			Message: err.Error(),
		})
	}

	// Get the user
	user := models.User{}
	result := database.DB.Model(&models.User{}).Where("email ? =", request.Email).First(&user)
	if result.Error != nil {
		message := "No account is associated with the email provided"

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = result.Error.Error()
		}

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{

			Message: message,
		})
	}

	// Check user's sign up provider
	if user.SignUpProvider != "local" {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{

			Message: "The sign up provider for this account does not support password reset",
		})
	}

	// Decrypt key
	str, err := facades.Crypt().DecryptString(request.Key)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Validate key, action, userId and timestamp
	split := strings.Split(str, "|")
	if len(split) != 3 || split[0] != "reset_password" || split[1] != user.ID || len(split[2]) == 0 || carbon.Now().DiffInMinutes(carbon.Parse(split[2])) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Invalid key provided",
		})
	}

	// Set new password
	passwordHash, err := utils.MakePasswordHash(request.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{

			Message: err.Error(),
		})
	}

	result = database.DB.Model(&user).Update("Password", passwordHash)

	// TODO: Delete all user tokens and log user out of all devices

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{

			Message: result.Error.Error(),
		})
	}

	return c.JSON(definitions.MessageResponse{
		Message: "Your password has been reset successfully.",
	})
}

func (a *AuthController) User(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)

	return c.JSON(models.UserToResource(user))
}

func (a *AuthController) VerifyAccount(c *fiber.Ctx) error {
	request := definitions.VerifyAccountRequest{}

	// Parse request
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Errors: errs,
		})
	}

	// Get user
	trx := database.DB.Begin()
	user := new(models.User)
	if err := trx.First(&user, "email = ?", request.Email).Error; err != nil {
		status := fiber.StatusBadRequest
		message := err.Error()

		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = fiber.StatusNotFound
			message = "No account is associated with the email provided"
		}

		return c.Status(status).JSON(definitions.MessageResponse{
			Message: message,
		})
	}

	// Check if user's sign up provider is local
	if user.SignUpProvider != "local" {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "The sign up provider for this account does not support account verification",
		})
	}

	// Check if user is already verified
	if user.EmailVerifiedAt != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Account already verified",
		})
	}

	// Decrypt key
	str, err := facades.Crypt().DecryptString(request.Key)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Validate key, action, userId and timestamp
	split := strings.Split(str, "|")
	if len(split) != 3 || split[0] != "account_verification" || split[1] != user.ID || len(split[2]) == 0 || carbon.Now().DiffInMinutes(carbon.Parse(split[2])) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Invalid key provided",
		})
	}

	// Update user
	if err := trx.Model(&user).Update("email_verified_at", time.Now()).Error; err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Send welcome email
	go jobs.NewMail().SendUserWelcomeMail(*user)

	// Commit transaction
	trx.Commit()

	return c.JSON(definitions.MessageResponse{
		Message: "Account verified successfully",
	})
}

func (a *AuthController) VerifyCode(c *fiber.Ctx) error {
	request := definitions.VerifyEmailRequest{}

	// Parse request
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Errors: errs,
		})
	}

	// Get user
	user := new(models.User)
	if err := database.DB.First(&user, "email = ?", request.Email).Error; err != nil {
		status := fiber.StatusBadRequest
		message := err.Error()

		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = fiber.StatusNotFound
			message = "No account is associated with the email provided"
		}

		return c.Status(status).JSON(definitions.MessageResponse{
			Message: message,
		})
	}

	// Check if code is valid
	code, err := config.RedisStore.Get(constants.USER_VERIFICATION_CODE_CACHE_KEY + user.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Invalid verification code",
		})
	}
	if request.Code != string(code) {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Invalid verification code",
		})
	}

	// Set encrypted code for password reset - action, userId, timestamp
	str, err := facades.Crypt().EncryptString(fmt.Sprintf("%s|%s|%s", request.Action, user.ID, time.Now().String()))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	config.RedisStore.Delete(constants.USER_VERIFICATION_CODE_CACHE_KEY + user.ID)

	return c.JSON(definitions.Map{
		"message": "Code verified successfully",
		"key":     str,
	})
}

func (a *AuthController) UpdateUserAvatar(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)

	// Parse request
	request := definitions.UpdateAvatarRequest{}
	avatar := ""
	key := ""
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Errors: errs,
		})
	}

	// Avatar key exists
	_, err := config.RedisStore.Get(request.AvatarKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Invalid avatar key",
		})
	}

	if request.Action == "remove" {
		// Delete avatar key from redis
		if err := config.RedisStore.Delete(request.AvatarKey); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: err.Error(),
			})
		}

		// Delete avatar from s3
		if err := utils.DeleteS3File(request.AvatarKey); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: err.Error(),
			})
		}
	}

	if request.Action == "update" {
		key = request.AvatarKey
		avatar, err = utils.GenerateS3FileURL(request.AvatarKey)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: err.Error(),
			})
		}
	}

	// Update user
	if err := database.DB.Model(&user).Update("avatar", key).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(definitions.Map{
		"avatarUrl": avatar,
		"message":   "Avatar updated successfully",
	})
}

func (a *AuthController) UpdateUserDetails(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	request := new(definitions.UpdateUserDetailsRequest)

	// Parse request
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Errors: errs,
		})
	}

	// Update user
	user.FirstName = request.FirstName
	user.LastName = request.LastName

	if len(request.PhoneNumber) == 10 && len(request.CountryCode) == 2 {
		// Parse phone number
		num, err := libphonenumber.Parse(request.PhoneNumber, strings.ToUpper(request.CountryCode))
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
				Errors: []definitions.ValidationErr{
					{Field: "phoneNumber", Reasons: []string{"Invalid phone number"}},
				},
			})
		}

		phoneNumber := libphonenumber.Format(num, libphonenumber.E164)

		user.PhoneNumber = phoneNumber
	}

	if err := database.DB.Model(&user).Updates(user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(models.UserToResource(user))
}

func (a *AuthController) SocialAuth(c *fiber.Ctx) error {
	// Parse request
	request := definitions.SocialAuthRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Errors: errs,
		})
	}

	// Get user
	if request.Provider == "google" {
		userInfo, err := facades.Socialite().Driver("google").UserFromToken(request.Token)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: err.Error(),
			})
		}

		// Check if user exists
		trx := database.DB.Begin()
		user := new(models.User)
		if err := trx.Model(&models.User{}).Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				trx.Rollback()

				return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
					Message: err.Error(),
				})
			}
		}

		if user.SignUpProvider == "local" {
			trx.Rollback()

			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: "An account with this email already exists with a different sign up provider",
			})
		}

		// Update or create user
		currentTime := time.Now()

		if user.EmailVerifiedAt == nil {
			user.EmailVerifiedAt = &currentTime
		}

		user.ExternalID = userInfo.ID
		user.Email = userInfo.Email
		user.FirstName = userInfo.FirstName
		user.LastName = userInfo.LastName
		user.SignUpProvider = "google"
		user.Channels = datatypes.JSON([]byte(`["web", "mobile"]`))

		if err := trx.Where("email = ?", userInfo.Email).Assign(user).FirstOrCreate(&user).Error; err != nil {
			trx.Rollback()

			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: err.Error(),
			})
		}

		// Generate token
		reference := nanoid.MustGenerate(constants.ALPHABETS, 32)
		accessToken, refreshToken, err := generateTokens(*user, reference)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: err.Error(),
			})
		}

		trx.Commit()

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"user": models.UserToResource(user),
			"tokens": fiber.Map{
				"access":  accessToken,
				"refresh": refreshToken,
			},
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
		Message: fmt.Sprintf("The selected provider %s is currently not supported", strings.ToUpper(request.Provider)),
	})
}

func (a *AuthController) Logout(c *fiber.Ctx) error {
	// Get token

	// Blacklist token

	return c.JSON(definitions.MessageResponse{

		Message: "Logout successful",
	})
}

func generateTokens(user models.User, reference string) (accessToken, refreshToken string, err error) {
	accessToken, err = models.GenerateUserToken(user, "access", reference)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = models.GenerateUserToken(user, "refresh", reference)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
