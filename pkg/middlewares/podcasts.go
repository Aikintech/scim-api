package middlewares

import (
	"encoding/json"
	// "strconv"
	// "strings"

	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/models"

	// "github.com/aikintech/scim/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// PodcastsCache is a middleware that caches podcasts
func PodcastsCache() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get podcasts from cache
		podcasts := make([]models.PodcastResource, 0)
		podcastsJson, err := config.RedisStore.Get(config.PODCASTS_CACHE_KEY)
		if err != nil {
			return c.Next()
		}

		err = json.Unmarshal(podcastsJson, &podcasts)
		if err != nil {
			return c.Next()
		}

		// podcasts = processPodcasts(c, podcasts)

		return c.JSON(definitions.DataResponse[[]models.PodcastResource]{
			Code: fiber.StatusOK,
			Data: podcasts,
		})
	}
}

// PodcastCache is a middleware that caches a particular podcast
func PodcastCache() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

// func processPodcasts(c *fiber.Ctx, podcasts []models.PodcastResource) []models.PodcastResource {
// 	limit, _ := strconv.Atoi(c.Query("limit", "10"))
// 	page, _ := strconv.Atoi(c.Query("page", "1"))
// 	sort := c.Query("sort", "newest")
// 	search := strings.Trim(c.Query("search", ""), " ")

// 	// Paginate podcasts

// 	return podcasts
// }
