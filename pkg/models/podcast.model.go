package models

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Podcast struct {
	ID          string `gorm:"primaryKey"`
	Author      string `gorm:"size:191;not null"`
	Title       string `gorm:"size:191;not null"`
	Description string `gorm:"type:LONGTEXT"`
	Duration    string `gorm:"size:191"`
	Meta        datatypes.JSON
	Status      bool `gorm:"type:tinyint"`
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Podcast) BeforeCreate(tx *gorm.DB) error {
	p.ID = ulid.Make().String()

	return nil
}
