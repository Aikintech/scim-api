package controllers

import (
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"time"
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
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Code:   fiber.StatusUnprocessableEntity,
			Errors: errs,
		})
	}
	//if len(request.ExcerptImageURL) > 0 {
	//	if !validation.IsValidFileKey(request.ExcerptImageURL) {
	//		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
	//			Code:   fiber.StatusUnprocessableEntity,
	//			Errors: []definitions.ValidationErr{{Field: "excerptImage", Reasons: []string{"Invalid excerpt image provided"}}},
	//		})
	//	}
	//}

	startDateTime, _ := time.Parse("2006-01-02 15:04:05", request.StartDateTime)
	endDateTime, _ := time.Parse("2006-01-02 15:04:05", request.EndDateTime)

	// Create event
	event := models.Event{
		Title:         request.Title,
		Description:   request.Description,
		Location:      request.Location,
		StartDateTime: startDateTime,
		EndDateTime:   &endDateTime,
	}

	return c.JSON(event)
}
