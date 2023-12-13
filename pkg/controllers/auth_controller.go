package controllers

import (
	"errors"
	"time"

	"github.com/aikintech/scim-api/pkg/config"
	"github.com/aikintech/scim-api/pkg/database"
	"github.com/aikintech/scim-api/pkg/jobs"

	"github.com/aikintech/scim-api/pkg/constants"

	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/utils"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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
	accessToken, err := models.GenerateUserToken(user, "access", reference)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}
	refreshToken, err := models.GenerateUserToken(user, "refresh", reference)
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
	user.Channels = datatypes.JSON([]byte(`["` + request.Channel + `"]`))
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

	// Set new password
	passwordHash, err := utils.MakePasswordHash(request.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{

			Message: err.Error(),
		})
	}

	result = database.DB.Model(&user).Update("Password", passwordHash)

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

	// Check if user is already verified
	if user.EmailVerifiedAt != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Account already verified",
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
