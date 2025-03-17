package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username string    `gorm:"uniqueIndex" json:"username"`
	Name     string    `json:"name"`
	Email    string    `gorm:"uniqueIndex" json:"email"`
	Password string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (user *User) BeforeCreate(tx *gorm.DB) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return nil
}