package controllers

import (
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/database"
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
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
func (evtCtrl *EventController) BackofficeGetEvents(c *fiber.Ctx) error {
	var total int64
	search := c.Query("search", "")

	// Query events
	events := make([]models.Event, 0)
	if err := database.DB.Model(&models.Event{}).Where("title LIKE ?", "%"+search+"%").Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}
	if err := database.DB.Model(&models.Event{}).Scopes(models.PaginationScope(c)).Where("title LIKE ?", "%"+search+"%").Find(&events).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(definitions.Map{
		"limit": c.QueryInt("limit", 10),
		"page":  c.QueryInt("page", 1),
		"total": total,
		"items": models.EventsToResourceCollection(events),
	})
}

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
		Title:           request.Title,
		Description:     request.Description,
		Location:        request.Location,
		StartDateTime:   startDateTime,
		EndDateTime:     &endDateTime,
		ExcerptImageURL: request.ExcerptImageURL,
	}
	if err := database.DB.Model(&models.Event{}).Create(&event).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(event.ToResource())
}

// BackofficeGetEvent
func (evtCtrl *EventController) BackofficeGetEvent(c *fiber.Ctx) error {
	eventId := c.Params("eventId")

	// Fetch event
	event := new(models.Event)
	if err := database.DB.Model(&models.Event{}).Where("id = ?", eventId).Find(&event).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(event.ToResource())
}

// BackofficeUpdateEvent
func (evtCtrl *EventController) BackofficeUpdateEvent(c *fiber.Ctx) error {
	eventId := c.Params("eventId")

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

	// Fetch event
	event := new(models.Event)
	if err := database.DB.Model(&models.Event{}).Where("id = ?", eventId).Find(&event).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	startDateTime, _ := time.Parse(constants.DATE_TIME_FORMAT, request.StartDateTime)
	endDateTime, _ := time.Parse(constants.DATE_TIME_FORMAT, request.EndDateTime)

	// Create event
	event.Title = request.Title
	event.Description = request.Description
	event.Location = request.Location
	event.StartDateTime = startDateTime
	event.EndDateTime = &endDateTime
	event.ExcerptImageURL = request.ExcerptImageURL

	database.DB.Save(&event)

	return c.Status(fiber.StatusCreated).JSON(event.ToResource())
}

func (evtCtrl *EventController) BackofficeDeleteEvent(c *fiber.Ctx) error {
	eventId := c.Params("eventId")

	// Fetch event
	event := new(models.Event)
	if err := database.DB.Model(&models.Event{}).Where("id = ?", eventId).Find(&event).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	if err := database.DB.Delete(&event).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(definitions.MessageResponse{
		Message: "Event deleted successfully",
	})
}
