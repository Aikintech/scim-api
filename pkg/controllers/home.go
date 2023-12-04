package controllers

import (
	"errors"
	"time"

	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HomeController struct{}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (homeCtrl *HomeController) ClientHome(c *fiber.Ctx) error {
	limit := 5
	upcomingEvents := []models.EventResource{}
	latestPodcasts := []models.PodcastResource{}
	latestPosts := []models.Post{}
	latestAnnouncements := []models.Post{}

	result := config.DB.Model(&models.Event{}).Where("start_date_time >= ?", time.Now()).Limit(limit).Find(&upcomingEvents)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Code:    fiber.StatusBadRequest,
				Message: result.Error.Error(),
			})
		}
	}

	// Latest podcasts (5)
	result = config.DB.Model(&models.Podcast{}).Order("published_at desc").Limit(limit).Find(&latestPodcasts)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Code:    fiber.StatusBadRequest,
				Message: result.Error.Error(),
			})
		}
	}

	// Latest blogPosts (5)
	result = config.DB.Model(&models.Post{}).Preload("User").Where("published = ?", true).Where("is_announcement = ?", false).Order("created_at desc").Limit(limit).Find(&latestPosts)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Code:    fiber.StatusBadRequest,
				Message: result.Error.Error(),
			})
		}
	}

	// Latest announcements (5)
	result = config.DB.Model(&models.Post{}).Preload("User").Where("published = ?", true).Where("is_announcement = ?", true).Order("created_at desc").Limit(limit).Find(&latestAnnouncements)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Code:    fiber.StatusBadRequest,
				Message: result.Error.Error(),
			})
		}
	}

	return c.JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"upcomingEvents":      upcomingEvents,
			"latestPodcasts":      latestPodcasts,
			"latestPosts":         latestPosts,
			"latestAnnouncements": latestAnnouncements,
		},
	})
}
