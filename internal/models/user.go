package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	Name               string     `json:"name" db:"name" validate:"required"`
	Email              string     `json:"email" db:"email" validate:"required,email"`
	Phone              *string    `json:"phone,omitempty" db:"phone"`
	Address            *string    `json:"address,omitempty" db:"address"`
	Role               string     `json:"role,omitempty" db:"role"`
	Status             string     `json:"status,omitempty" db:"status"`
	SubscriptionStatus string     `json:"subscription_status,omitempty" db:"subscription_status"`
	Password           string     `json:"password,omitempty" db:"password" validate:"required"`
	CreatedAt          time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
