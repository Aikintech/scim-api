package models

import (
	"fmt"
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/utils"
	nanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Testimony struct {
	ID                 string  `gorm:"primaryKey;size:40"`
	YouTubeReferenceID *string `gorm:"size:40;column:yt_reference_id"`
	TikTokReferenceID  *string `gorm:"size:40;column:tk_reference_id"`
	YouTubeURL         *string
	TikTokURL          *string
	Title              string    `gorm:"size:191;not null"`
	FileURL            string    `gorm:"size:191;not null"`
	ThumbnailURL       string    `gorm:"size:191"`
	Body               string    `gorm:"not null"`
	Published          bool      `gorm:"not null;default:false"`
	CreatedAt          time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime;not null"`
}

type TestimonyResource struct {
	ID           string    `json:"id"`
	YouTubeURL   *string   `json:"youtubeUrl,omitempty"`
	TikTokURL    *string   `json:"tiktokUrl,omitempty"`
	Title        string    `json:"title"`
	FileURL      string    `json:"fileUrl"`
	FileKey      string    `json:"fileKey"`
	ThumbnailURL *string   `json:"thumbnailUrl,omitempty"`
	ThumbnailKey *string   `json:"thumbnailKey,omitempty"`
	Body         string    `json:"body"`
	Published    bool      `json:"published"`
	CreatedAt    time.Time `json:"createAt"`
}

func (t *Testimony) BeforeCreate(tx *gorm.DB) error {
	t.ID = nanoid.MustGenerate(constants.ALPHABETS, 32)

	return nil
}

func TestimonyToResource(t *Testimony) *TestimonyResource {
	// Generate fileURL
	fileURL, err := utils.GenerateS3FileURL(t.FileURL)
	if err != nil {
		fmt.Println("Error generating avatar url", err.Error())
	}

	return &TestimonyResource{
		ID:         t.ID,
		YouTubeURL: t.YouTubeURL,
		TikTokURL:  t.TikTokURL,
		Title:      t.Title,
		FileURL:    fileURL,
		FileKey:    t.FileURL,
		Body:       t.Body,
		Published:  t.Published,
		CreatedAt:  t.CreatedAt,
	}
}

func TestimoniesToResourceCollection(testimonies []*Testimony) []*TestimonyResource {
	results := make([]*TestimonyResource, 0)

	for _, t := range testimonies {
		results = append(results, TestimonyToResource(t))
	}

	return results
}
