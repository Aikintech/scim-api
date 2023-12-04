package controllers

import "github.com/gofiber/fiber/v2"

type EventController struct{}

func NewEventController() *EventController {
	return &EventController{}
}

func (evtCtrl *EventController) GetEvents(c *fiber.Ctx) error {
	return c.SendString("GetEvents")
}

// Backoffice
func (evtCtrl *EventController) BackofficeCreateEvent(c *fiber.Ctx) error {
	return c.SendString("CreateEvent")
}
