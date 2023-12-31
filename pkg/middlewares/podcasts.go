package middlewares

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/aikintech/scim-api/pkg/constants"

	"github.com/aikintech/scim-api/pkg/config"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func AllPodcastsCache() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Query params
		orderBy := c.Query("sort", "newest")

		// Get podcasts from cache
		podcasts := make([]models.PodcastResource, 0)
		podcastsJson, err := config.RedisStore.Get(constants.PODCASTS_CACHE_KEY)
		if err != nil {
			return c.Next()
		}

		err = json.Unmarshal(podcastsJson, &podcasts)
		if err != nil {
			return c.Next()
		}

		// Sort podcasts
		if orderBy == "newest" {
			sort.Slice(podcasts, func(i, j int) bool {
				return podcasts[i].PublishedAt.String() > podcasts[j].PublishedAt.String()
			})
		} else {
			sort.Slice(podcasts, func(i, j int) bool {
				return podcasts[i].PublishedAt.String() < podcasts[j].PublishedAt.String()
			})
		}

		fmt.Println("Podcasts from cache")

		return c.JSON(podcasts)
	}
}

// PodcastsCache is a middleware that caches podcasts
func DBPodcastsCache() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get podcasts from cache
		podcasts := make([]models.PodcastResource, 0)
		podcastsJson, err := config.RedisStore.Get(constants.PODCASTS_CACHE_KEY)
		if err != nil {
			return c.Next()
		}

		err = json.Unmarshal(podcastsJson, &podcasts)
		if err != nil {
			return c.Next()
		}

		// podcasts = processPodcasts(c, podcasts)

		return c.JSON(podcasts)
	}
}

// PodcastCache is a middleware that caches a particular podcast
func PodcastByIdCache() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		podcastId := c.Params("podcastId", "")

		// Get podcasts from cache
		podcasts := make([]models.PodcastResource, 0)
		podcastsJson, err := config.RedisStore.Get(constants.PODCASTS_CACHE_KEY)
		if err != nil {
			return c.Next()
		}

		err = json.Unmarshal(podcastsJson, &podcasts)
		if err != nil {
			return c.Next()
		}

		// Find podcast
		podcast := models.PodcastResource{}
		for _, p := range podcasts {
			if p.ID == podcastId {
				podcast = p
				break
			}
		}

		return c.JSON(podcast)
	}
}
