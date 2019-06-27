package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Secret model contains all information about itself.
type Secret struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`

	// Unique hash to identify the secrets
	Hash string `json:"hash" db:"hash" gorm:"unique_index:idx_hash"`

	// The secret itself
	SecretText string `json:"secretText"`

	// The date and time of the creation
	CreatedAt time.Time `json:"createdAt"  time_format:"2006-01-02 15:04:05"`

	// The secret cannot be reached after this time
	ExpiresAt time.Time `json:"expiresAt"  time_format:"2006-01-02 15:04:05"`

	// How many times the secret can be viewed
	RemainingViews int32 `json:"remainingViews"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (s *Secret) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid.String())
}
