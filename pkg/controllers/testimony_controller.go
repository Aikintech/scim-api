package controllers

import (
	"fmt"
	"mime/multipart"
	"slices"
	"strconv"
	"strings"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/database"
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/facades"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/utils"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TestimonyController struct {
	query *gorm.DB
}

type uploadAndCreateType struct {
	request  definitions.TestimonyRequest
	file     multipart.File
	filename string
	user     *models.User
}

func NewTestimonyController() *TestimonyController {
	q := database.DB.Model(&models.Testimony{}).
		Select("testimonies.*, COUNT(DISTINCT likes.id) AS likes_count, COUNT(DISTINCT comments.id) AS comments_count").
		Joins("JOIN likes ON likes.likeable_id = testimonies.id AND likes.likeable_type = 'testimonies'").
		Joins("JOIN comments ON comments.commentable_id = testimonies.id AND comments.commentable_type = 'testimonies'").
		Group("testimonies.id")

	return &TestimonyController{
		query: q,
	}
}

func (t *TestimonyController) GetAllTestimonies() {}

func (t *TestimonyController) GetTestimonies() {}

func (t *TestimonyController) GetTestimony() {}

// Backoffice routes
func (t *TestimonyController) BackofficeGetTestimonies(c *fiber.Ctx) error {
	var total int64
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)
	search := strings.TrimSpace(c.Query("search", ""))
	query := t.query.Where("testimonies.title LIKE ? OR testimonies.body LIKE ?", "%"+search+"%", "%"+search+"%")

	// Query total
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Query results
	testimonies := make([]*models.Testimony, 0)
	if err := query.Scopes(models.PaginationScope(c)).Find(&testimonies).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(definitions.Map{
		"limit": limit,
		"page":  page,
		"total": total,
		"items": models.TestimoniesToResourceCollection(testimonies),
	})
}

func (t *TestimonyController) BackofficeCreateTestimony(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	published, _ := strconv.ParseBool(c.FormValue("published"))
	request := definitions.TestimonyRequest{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Published:   published,
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Errors: errs,
		})
	}

	// Get file from request
	requestFile, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "File upload failed or file not present",
		})
	}

	// Validate mime type
	mime := utils.GetMimeExtension(requestFile.Header["Content-Type"][0])
	if !slices.Contains([]string{"mov", "mp4"}, mime) {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Invalid file provided. Expected: mov, mp4",
		})
	}

	// Validate file size
	fileSize := requestFile.Size / 1024 / 1024
	if fileSize > 60 {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "file size must be less than 60MB",
		})
	}

	file, err := requestFile.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Generate unique filename
	filename := strings.ToLower(fmt.Sprintf("%s/%s.%s", "testimony", utils.GenerateRandomString(30), mime))

	go t.storeTestimony(uploadAndCreateType{
		request:  request,
		file:     file,
		filename: filename,
		user:     user,
	})

	return c.JSON(definitions.MessageResponse{
		Message: "Upload of testimony has started. You will be notified when it is finished.",
	})
}

func (t *TestimonyController) storeTestimony(data uploadAndCreateType) {
	fmt.Println("Store testimony started...")
	result, err := utils.UploadFileS3(data.file, data.filename)

	if err != nil {
		fmt.Printf("Store testimony ended with error: %s", err.Error())

		pData := facades.PusherData{"notificationType": "error", "message": err.Error()}
		facades.Pusher().Trigger(data.user.ID, constants.NOTIFICATION_EVENT, pData)
	} else {
		trx := database.DB.Begin()
		testimony := models.Testimony{
			Title:     data.request.Title,
			Body:      data.request.Description,
			Published: data.request.Published,
			FileURL:   result.Key,
		}

		err := trx.Model(&models.Testimony{}).Create(&testimony).Error
		if err != nil {
			fmt.Printf("Store testimony ended with error: %s", err.Error())

			trx.Rollback()

			e := utils.DeleteS3File(result.Key)
			if e != nil {
				pData := facades.PusherData{"notificationType": "error", "message": e.Error()}
				facades.Pusher().Trigger(data.user.ID, constants.NOTIFICATION_EVENT, pData)
			} else {
				pData := facades.PusherData{"notificationType": "error", "message": err.Error()}
				facades.Pusher().Trigger(data.user.ID, constants.NOTIFICATION_EVENT, pData)
			}
		} else {
			trx.Commit()

			pData := facades.PusherData{
				"notificationType": "success",
				"message":          fmt.Sprintf("Testimony %s has been uploaded successfully.", testimony.Title),
				"module":           "testimonies",
			}
			facades.Pusher().Trigger(data.user.ID, constants.NOTIFICATION_EVENT, pData)

			fmt.Println("Store testimony ended with success")
		}
	}
}
