package controllers

import (
	"github.com/aikintech/scim-api/pkg/jobs"
	"github.com/gofiber/fiber/v2"
)

type MiscController struct{}

func NewMiscController() *MiscController {
	return &MiscController{}
}

func (miscCtrl *MiscController) HealthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func (miscCtrl *MiscController) BackupDatabase(c *fiber.Ctx) error {
	go jobs.BackupDatabase()

	return c.SendString("Backup initiated")
}

func (miscCtrl *MiscController) SeedPodcasts(c *fiber.Ctx) error {
	go jobs.SeedPodcasts()

	return c.SendString("Podcasts seeding initiated")
}
