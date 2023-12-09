package models

type PodcastPlaylist struct {
	PlaylistID string `gorm:"column:playlist_id"`
	PodcastID  string `gorm:"column:podcast_id"`
	Podcast    Podcast
	Playlist   Playlist
}

func (PodcastPlaylist) TableName() string {
	return "podcast_playlist"
}
