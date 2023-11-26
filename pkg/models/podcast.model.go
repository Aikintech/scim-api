package models

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Podcast struct {
	ID          string     `gorm:"primaryKey"`
	GUID        string     `gorm:"size:191;not null"`
	Author      string     `gorm:"size:191;not null"`
	Title       string     `gorm:"size:191;not null"`
	SubTitle    string     `gorm:"size:191;not null"`
	Summary     string     `gorm:"size:191;not null"`
	Description string     `gorm:"type:LONGTEXT"`
	Duration    string     `gorm:"size:191;not null"`
	Image       string     `gorm:"size:191;not null"`
	Url         string     `gorm:"type:TEXT;not null"`
	Published   bool       `gorm:"type:tinyint; not null;default:true"`
	PublishedAt *time.Time `gorm:"not null"`
	Meta        datatypes.JSON
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Podcast) BeforeCreate(tx *gorm.DB) error {
	p.ID = ulid.Make().String()

	return nil
}
