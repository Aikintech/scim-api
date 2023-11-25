package definitions

import "time"

type PodcastFeedItem struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Link            string    `json:"link"`
	Links           []string  `json:"links"`
	Published       string    `json:"published"`
	PublishedParsed time.Time `json:"publishedParsed"`
	Author          struct {
		Name string `json:"name"`
	} `json:"author"`
	Authors []struct {
		Name string `json:"name"`
	} `json:"authors"`
	GUID  string `json:"guid"`
	Image struct {
		URL string `json:"url"`
	} `json:"image"`
	Enclosures []struct {
		URL    string `json:"url"`
		Length string `json:"length"`
		Type   string `json:"type"`
	} `json:"enclosures"`
	ItunesExt struct {
		Author      string `json:"author"`
		Duration    string `json:"duration"`
		Subtitle    string `json:"subtitle"`
		Summary     string `json:"summary"`
		Image       string `json:"image"`
		EpisodeType string `json:"episodeType"`
	} `json:"itunesExt"`
	Extensions struct {
		Itunes struct {
			Author []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
				Attrs struct {
				} `json:"attrs"`
				Children struct {
				} `json:"children"`
			} `json:"author"`
			Duration []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
				Attrs struct {
				} `json:"attrs"`
				Children struct {
				} `json:"children"`
			} `json:"duration"`
			EpisodeType []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
				Attrs struct {
				} `json:"attrs"`
				Children struct {
				} `json:"children"`
			} `json:"episodeType"`
			Image []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
				Attrs struct {
					Href string `json:"href"`
				} `json:"attrs"`
				Children struct {
				} `json:"children"`
			} `json:"image"`
			Subtitle []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
				Attrs struct {
				} `json:"attrs"`
				Children struct {
				} `json:"children"`
			} `json:"subtitle"`
			Summary []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
				Attrs struct {
				} `json:"attrs"`
				Children struct {
				} `json:"children"`
			} `json:"summary"`
			Title []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
				Attrs struct {
				} `json:"attrs"`
				Children struct {
				} `json:"children"`
			} `json:"title"`
		} `json:"itunes"`
	} `json:"extensions"`
}
