package jobs

import (
	"fmt"
	"github.com/aikintech/scim/pkg/config"
	"github.com/mmcdole/gofeed"
)

func SeedPodcasts() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(config.PodcastUrl)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(feed.FeedVersion)
}
