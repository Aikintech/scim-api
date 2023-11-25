package jobs

import (
	"github.com/aikintech/scim/pkg/config"
	"github.com/mmcdole/gofeed"
	"os"
)

func SeedPodcasts() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(config.PodcastUrl)

	if err != nil {
		panic(err.Error())
	}

	if err := os.WriteFile("podcast.json", []byte(feed.String()), 0777); err != nil {
		panic(err.Error())
	}
}
