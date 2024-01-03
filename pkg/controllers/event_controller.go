package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/database"
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/utils"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type EventController struct{}

func NewEventController() *EventController {
	return &EventController{}
}

func (evtCtrl *EventController) GetEvents(c *fiber.Ctx) error {
	var total int64
	search := c.Query("search", "")

	// Query events
	events := make([]*models.Event, 0)
	query := database.DB.Model(&models.Event{}).
		Where("title LIKE ?", "%"+search+"%").
		Where("start_date_time >= DATE(?)", time.Now())

	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}
	if err := query.Scopes(models.PaginationScope(c)).Preload("Users").Find(&events).Error; err != nil {
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

func (evtCtrl *EventController) GetEvent(c *fiber.Ctx) error {
	eventId := c.Params("eventId")

	// Fetch event
	event := new(models.Event)
	if err := database.DB.Model(&models.Event{}).Preload("Users").Where("id = ?", eventId).Find(&event).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(event.ToResource())
}

func (evtCtrl *EventController) MyCalendarEvents(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	events := make([]*models.Event, 0)

	if err := database.DB.
		Model(&models.Event{}).
		Select("events.*").
		Joins("JOIN user_event ON events.id = user_event.event_id").
		Where("user_event.user_id = ?", user.ID).
		Group("events.id").Find(&events).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: err.Error(),
			})
		}
	}

	return c.JSON(models.EventsToResourceCollection(events))
}

func (evtCtrl *EventController) SyncEventToCalendar(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	eventId := c.Params("eventId")
	trx := database.DB.Begin()
	action := "added"

	// Find event
	event := new(models.Event)
	if err := trx.Model(&models.Event{}).Where("id = ?", eventId).First(&event).Error; err != nil {
		status := fiber.StatusNotFound
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			status = fiber.StatusBadRequest
		}
		return c.Status(status).JSON(definitions.MessageResponse{Message: err.Error()})
	}

	// Find user event
	userEvent := new(models.UserEvent)
	if err := trx.Model(&models.UserEvent{}).Where("user_id = ? AND event_id = ?", user.ID, eventId).First(&userEvent).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{Message: err.Error()})
		}
	}

	// Delete or create user_event
	if userEvent.UserID != "" && userEvent.EventID != "" {
		if err := trx.Where("user_id = ? AND event_id = ?", user.ID, eventId).Delete(&userEvent).Error; err != nil {
			trx.Rollback()

			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{Message: err.Error()})
		}

		action = "removed"
	} else {
		if err := trx.Model(&models.UserEvent{}).Create(models.UserEvent{UserID: user.ID, EventID: eventId}).Error; err != nil {
			trx.Rollback()

			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{Message: err.Error()})
		}
	}

	trx.Commit()

	return c.JSON(definitions.MessageResponse{
		Message: fmt.Sprintf("Event %s to calendar successfully", action),
	})
}

/*** Backoffice handlers ***/
func (evtCtrl *EventController) BackofficeGetEvents(c *fiber.Ctx) error {
	var total int64
	search := c.Query("search", "")

	// Query events
	events := make([]*models.Event, 0)
	query := database.DB.Model(&models.Event{}).Where("title LIKE ?", "%"+search+"%")
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}
	if err := query.Scopes(models.PaginationScope(c)).Order("created_at DESC").Find(&events).Error; err != nil {
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
		Published:       request.Published,
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
	event.Published = request.Published
	event.StartDateTime = startDateTime
	event.EndDateTime = &endDateTime
	event.ExcerptImageURL = request.ExcerptImageURL

	database.DB.Save(&event)

	return c.Status(fiber.StatusCreated).JSON(event.ToResource())
}

func (evtCtrl *EventController) BackofficeDeleteEvent(c *fiber.Ctx) error {
	eventId := c.Params("eventId")

	// Fetch event
	trx := database.DB.Begin()
	event := new(models.Event)
	if err := trx.Model(&models.Event{}).Where("id = ?", eventId).Find(&event).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	if err := trx.Delete(&models.Event{}, "id = ?", eventId).Error; err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Delete file from storage
	go func() {
		if err := utils.DeleteS3File(event.ExcerptImageURL); err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("File deletion done...")
	}()

	trx.Commit()

	return c.JSON(definitions.MessageResponse{
		Message: "Event deleted successfully",
	})
}
