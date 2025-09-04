package entities

import (
	"time"
)

type ActivityType struct {
	BaseEntity
	Name              string  `json:"name" gorm:"not null;uniqueIndex"`
	CaloriesPerMinute float64 `json:"caloriesPerMinute" gorm:"not null"`
	Description       *string `json:"description"`
}

type Activity struct {
	BaseEntity
	UserID            uint         `json:"userId" gorm:"not null"`
	ActivityTypeID    uint         `json:"activityTypeId" gorm:"not null"`
	DurationInMinutes int          `json:"durationInMinutes" gorm:"not null"`
	CaloriesBurned    float64      `json:"caloriesBurned" gorm:"not null"`
	DoneAt            time.Time    `json:"doneAt" gorm:"not null"`
	
	User              User         `json:"user" gorm:"foreignKey:UserID"`
	ActivityType      ActivityType `json:"activityType" gorm:"foreignKey:ActivityTypeID"`
}

type CreateActivityRequest struct {
	ActivityType      string `json:"activityType" binding:"required"`
	DoneAt            string `json:"doneAt" binding:"required"`
	DurationInMinutes int    `json:"durationInMinutes" binding:"required,min=1"`
}

type ActivityResponse struct {
	ActivityID        string  `json:"activityId"`
	ActivityType      string  `json:"activityType"`
	DoneAt            string  `json:"doneAt"`
	DurationInMinutes int     `json:"durationInMinutes"`
	CaloriesBurned    float64 `json:"caloriesBurned"`
	CreatedAt         string  `json:"createdAt"`
}

type ActivityFilter struct {
	UserID       uint
	ActivityType *string
	StartDate    *time.Time
	EndDate      *time.Time
}