package models

import (
	"time"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/oklog/ulid/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Podcast struct {
	ID          string     `gorm:"primaryKey;size:40"`
	GUID        string     `gorm:"size:191;not null"`
	Author      string     `gorm:"size:191;not null"`
	Title       string     `gorm:"not null"`
	SubTitle    string     `gorm:"not null"`
	Summary     string     `gorm:"not null"`
	Description string     `gorm:"type:TEXT"`
	Duration    string     `gorm:"size:191;not null"`
	ImageURL    string     `gorm:"not null"`
	AudioURL    string     `gorm:"type:TEXT;not null"`
	Published   bool       `gorm:"not null;default:true"`
	PublishedAt *time.Time `gorm:"not null"`
	ShortURL    string     `gorm:"size:191"`
	Meta        datatypes.JSON
	CreatedAt   time.Time   `gorm:"not null"`
	UpdatedAt   time.Time   `gorm:"not null"`
	Playlists   []*Playlist `gorm:"many2many:podcast_playlist"`
	Comments    []Comment   `gorm:"polymorphic:Commentable"`
	Likes       []Like      `gorm:"polymorphic:Likeable"`
}

func (p *Podcast) BeforeCreate(tx *gorm.DB) error {
	p.ID = ulid.Make().String()

	if len(p.Summary) > 0 {
		p.Summary = strip.StripTags(p.Summary)
	}

	if len(p.Description) > 0 {
		p.Description = strip.StripTags(p.Description)
	}

	return nil
}

func (p *Podcast) GetCommentableType() string {
	return "podcasts"
}

type PodcastResource struct {
	ID          string    `json:"id"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Duration    string    `json:"duration"`
	ImageURL    string    `json:"imageUrl"`
	AudioURL    string    `json:"audioUrl"`
	Published   bool      `json:"published"`
	PublishedAt time.Time `json:"publishedAt"`
}