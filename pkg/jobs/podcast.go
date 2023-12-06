package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aikintech/scim-api/pkg/constants"
	"time"

	"github.com/aikintech/scim-api/pkg/config"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func SeedPodcasts() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(constants.PODCAST_URL)

	if err != nil {
		panic(err.Error())
	}

	podcasts := make([]models.Podcast, 0)

	// Loop through podcasts
	for _, items := range lo.Chunk(feed.Items, 100) {
		for _, item := range items {
			var podcast = models.Podcast{
				GUID:        item.GUID,
				Author:      item.ITunesExt.Author,
				Title:       item.Title,
				SubTitle:    item.ITunesExt.Subtitle,
				Summary:     item.ITunesExt.Summary,
				Description: item.Description,
				Duration:    item.ITunesExt.Duration,
				ImageURL:    item.Image.URL,
				AudioURL:    item.Enclosures[0].URL,
				Published:   true,
				PublishedAt: item.PublishedParsed,
			}

			updateOrCreate(&podcast)
			podcasts = append(podcasts, podcast)
		}
	}

	// Convert podcasts to JSON
	podcastsJson, err := json.Marshal(podcasts)
	if err != nil {
		panic(err.Error())
	}

	// Cache podcasts for 24 hours
	config.RedisStore.Set(constants.PODCASTS_CACHE_KEY, podcastsJson, time.Hour*24)

	fmt.Println("Podcasts seeded successfully")
}

func updateOrCreate(podcast *models.Podcast) {
	result := config.DB.First(&podcast, "guid = ?", podcast.GUID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result := config.DB.Debug().Create(&podcast)

			if result.Error != nil {
				fmt.Println(result.Error.Error())
			}
		}
	}
}
