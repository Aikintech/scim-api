package controllers

import (
	"fmt"
	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/validation"
	"strings"

	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type FileController struct {
}

func NewFileController() *FileController {
	return &FileController{}
}

func (fileCtrl *FileController) UploadFile(c *fiber.Ctx) error {
	// Get file from request
	requestFile, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "File upload failed",
		})
	}

	// Validate upload type
	uploadType := c.FormValue("uploadType")
	if !lo.Contains(constants.UPLOAD_TYPES, uploadType) {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid upload type",
		})
	}

	// Validate mime type
	mime := utils.GetMimeExtension(requestFile.Header["Content-Type"][0])
	if !lo.Contains(constants.ALLOWED_MIME_TYPES, mime) {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid file type",
		})
	}

	// Validate file size
	fileSize := requestFile.Size / 1024 / 1024
	fileErrMsg := ""
	if uploadType == "testimony" && fileSize > 60 {
		fileErrMsg = "File size must be less than 60MB"
	}
	if uploadType != "testimony" && fileSize > 1 {
		fileErrMsg = "File size must be less than 1MB"
	}
	if len(fileErrMsg) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: fileErrMsg,
		})
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s/%s.%s", strings.ToUpper(uploadType), utils.GenerateRandomString(30), mime)

	// Upload to s3
	result, err := utils.UploadFileS3(requestFile, filename)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.JSON(definitions.Map{
		"code": fiber.StatusOK,
		"data": result,
	})
}

func (fileCtrl *FileController) GetFileURL(c *fiber.Ctx) error {
	key := c.Query("key", "")

	// Validate request
	if validation.IsValidFileKey(key) {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid file key",
		})
	}

	// Generate file URL
	location, err := utils.GenerateS3FileURL(key)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.JSON(definitions.Map{
		"code": fiber.StatusOK,
		"data": map[string]string{
			"key": key,
			"url": location,
		},
	})
}
