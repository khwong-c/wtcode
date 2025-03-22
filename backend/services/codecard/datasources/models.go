package datasources

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CodeCard struct {
	ID        uuid.UUID `gorm:"type:uuid"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Title       string
	Code        string
	Language    string
	Description *string
	Example     *string

	Author      *string
	AuthorEmail *string
	AuthorURL   *string

	Verified bool
}
