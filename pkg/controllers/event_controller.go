package controllers

import (
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/database"
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/utils"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type EventController struct{}

func NewEventController() *EventController {
	return &EventController{}
}

func (evtCtrl *EventController) GetEvents(c *fiber.Ctx) error {
	return c.SendString("GetEvents")
}

/*** Backoffice handlers ***/
// BackofficeCreateEvent - add a new event
func (evtCtrl *EventController) BackofficeCreateEvent(c *fiber.Ctx) error {
	// Parse request
	request := new(definitions.StoreEventRequest)
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

	startDateTime, _ := time.Parse(constants.DATE_TIME_FORMAT, request.StartDateTime)
	endDateTime, _ := time.Parse(constants.DATE_TIME_FORMAT, request.EndDateTime)

	// Create event
	event := models.Event{
		Title:         request.Title,
		Description:   request.Description,
		Location:      request.Location,
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
	}
	if err := database.DB.Model(&models.Event{}).Create(&event).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	excerptImage, _ := utils.GenerateS3FileURL(request.ExcerptImageURL)
	return c.Status(fiber.StatusCreated).JSON(definitions.DataResponse[*models.EventResource]{
		Code: fiber.StatusCreated,
		Data: &models.EventResource{
			ID:              event.ID,
			Title:           event.Title,
			Description:     event.Description,
			ExcerptImageURL: &excerptImage,
			Location:        event.Location,
			StartDateTime:   event.StartDateTime,
			EndDateTime:     &event.EndDateTime,
			CreatedAt:       event.CreatedAt,
		},
	})
}
