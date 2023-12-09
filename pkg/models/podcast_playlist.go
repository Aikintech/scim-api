package models

type PodcastPlaylist struct {
	PlaylistID string `gorm:"column:playlist_id;index"`
	PodcastID  string `gorm:"column:podcast_id;index"`
	Podcast    Podcast
	Playlist   Playlist
}

func (PodcastPlaylist) TableName() string {
	return "podcast_playlist"
}
