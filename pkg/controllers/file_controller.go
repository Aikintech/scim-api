package controllers

import (
	"fmt"
	"slices"
	"strings"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/validation"

	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type FileController struct {
}

func NewFileController() *FileController {
	return &FileController{}
}

func (fileCtrl *FileController) UploadFile(c *fiber.Ctx) error {
	path := c.Path()
	uploadType := c.FormValue("uploadType")
	hasAvatarPrefix := strings.HasPrefix(strings.ToLower(uploadType), "avatar")
	hasBackofficePrefix := strings.HasPrefix(path, "/backoffice")

	// Validate upload type
	if !slices.Contains(constants.UPLOAD_TYPES, uploadType) {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Invalid upload type",
		})
	}

	// Only backoffice can upload EXCERPT and TESTIMONY files
	if !hasBackofficePrefix && !hasAvatarPrefix {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Invalid upload type",
		})
	}

	// Get file from request
	requestFile, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "File upload failed",
		})
	}

	// Validate mime type
	mime := utils.GetMimeExtension(requestFile.Header["Content-Type"][0])
	if !slices.Contains(constants.ALLOWED_MIME_TYPES, mime) {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
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
			Message: fileErrMsg,
		})
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s/%s.%s", strings.ToUpper(uploadType), utils.GenerateRandomString(30), mime)

	// Upload to s3
	result, err := utils.UploadFileS3(requestFile, filename)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(result)
}

func (fileCtrl *FileController) GetFileURL(c *fiber.Ctx) error {
	key := c.Params("fileKey", "")

	// Validate request
	if !validation.IsValidFileKey(key) {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Invalid file key",
		})
	}

	// Generate file URL
	location, err := utils.GenerateS3FileURL(key)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(definitions.Map{
		"key": key,
		"url": location,
	})
}

func (fileCtrl *FileController) DeleteFile(c *fiber.Ctx) error {
	path := c.Path()
	key := strings.TrimSpace(c.Params("fileKey", ""))
	hasAvatarPrefix := strings.HasPrefix(strings.ToLower(key), "avatar")
	hasBackofficePrefix := strings.HasPrefix(path, "/backoffice")

	// Validate request
	if !validation.IsValidFileKey(key) {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Invalid file key",
		})
	}

	// Only backoffice can delete TESTIMONY and EXCERPT files
	if !hasBackofficePrefix && !hasAvatarPrefix {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Invalid file key",
		})
	}

	// Delete file from s3
	if err := utils.DeleteS3File(key); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(definitions.MessageResponse{
		Message: "File deleted successfully",
	})
}
