package entities

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Email      string    `json:"email" gorm:"uniqueIndex;not null"`
	Name       *string   `json:"name" gorm:"type:varchar(255)"`
	Preference *string   `json:"preference" gorm:"type:varchar(255)"`
	WeightUnit *string   `json:"weightUnit" gorm:"type:varchar(10)"`
	HeightUnit *string   `json:"heightUnit" gorm:"type:varchar(10)"`
	Weight     *float64  `json:"weight" gorm:"type:decimal(5,2)"`
	Height     *float64  `json:"height" gorm:"type:decimal(5,2)"`
	ImageURI   *string   `json:"imageUri" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Email      string   `json:"email" binding:"required,email"`
	Name       *string  `json:"name"`
	Preference *string  `json:"preference"`
	WeightUnit *string  `json:"weightUnit"`
	HeightUnit *string  `json:"heightUnit"`
	Weight     *float64 `json:"weight"`
	Height     *float64 `json:"height"`
	ImageURI   *string  `json:"imageUri"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Email      *string  `json:"email,omitempty" binding:"omitempty,email"`
	Name       *string  `json:"name,omitempty"`
	Preference *string  `json:"preference,omitempty"`
	WeightUnit *string  `json:"weightUnit,omitempty"`
	HeightUnit *string  `json:"heightUnit,omitempty"`
	Weight     *float64 `json:"weight,omitempty"`
	Height     *float64 `json:"height,omitempty"`
	ImageURI   *string  `json:"imageUri,omitempty"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID         uint     `json:"id"`
	Email      string   `json:"email"`
	Name       *string  `json:"name"`
	Preference *string  `json:"preference"`
	WeightUnit *string  `json:"weightUnit"`
	HeightUnit *string  `json:"heightUnit"`
	Weight     *float64 `json:"weight"`
	Height     *float64 `json:"height"`
	ImageURI   *string  `json:"imageUri"`
}
