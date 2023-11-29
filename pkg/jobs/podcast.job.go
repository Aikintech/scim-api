package jobs

import (
	"errors"
	"fmt"

	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/models"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func SeedPodcasts() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(config.PodcastUrl)

	if err != nil {
		panic(err.Error())
	}

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
				Image:       item.Image.URL,
				Url:         item.Enclosures[0].URL,
				Published:   true,
				PublishedAt: item.PublishedParsed,
			}

			updateOrCreate(&podcast)
		}
	}
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
