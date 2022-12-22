package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.Id == "" {
		b.Id = uuid.New().String()
	}
	return
}
