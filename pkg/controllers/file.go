package controllers

import (
	"fmt"
	"strings"

	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/utils"
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
	if !lo.Contains(config.UPLOAD_TYPES, uploadType) {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid upload type",
		})
	}

	// Validate mime type
	mime := utils.GetMimeExtension(requestFile.Header["Content-Type"][0])
	if !lo.Contains(config.ALLOWED_MIME_TYPES, mime) {
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
	if !lo.SomeBy(config.UPLOAD_TYPES, func(item string) bool {
		uploadType := strings.ToUpper(item)

		return strings.Contains(key, uploadType) // Key contains upload type
	}) {
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
