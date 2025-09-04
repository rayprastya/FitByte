package repositories

import (
	"context"
	"fitbyte/internal/entities"
	"time"

	"gorm.io/gorm"
)

type (
	ActivityRepository interface {
		GetActivitiesWithFilters(ctx context.Context, filter ActivityFilter) ([]entities.Activity, int64, error)
		GetActivity(ctx context.Context, filter ActivityFilter) (*entities.Activity, error)
		CreateActivity(ctx context.Context, activity entities.Activity) (*entities.Activity, error)
		UpdateActivity(ctx context.Context, activity entities.Activity) (*entities.Activity, error)
		DeleteActivity(ctx context.Context, id uint) error
		GetActivityType(ctx context.Context, name string) (*entities.ActivityType, error)
		GetActivityTypes(ctx context.Context) ([]entities.ActivityType, error)
		CreateActivityType(ctx context.Context, activityType entities.ActivityType) (*entities.ActivityType, error)
	}

	activityRepo struct {
		db *gorm.DB
	}

	ActivityRepoParam struct {
		DB *gorm.DB
	}

	ActivityFilter struct {
		UserID             *uint
		ActivityType       *string
		Limit              int
		Offset             int
		DoneAtFrom         *time.Time
		DoneAtTo           *time.Time
		CaloriesBurnedMin  *float64
		CaloriesBurnedMax  *float64
		ID                 *uint
	}
)

func NewActivityRepository(param ActivityRepoParam) ActivityRepository {
	return &activityRepo{
		db: param.DB,
	}
}

func (repo *activityRepo) GetActivitiesWithFilters(ctx context.Context, filter ActivityFilter) ([]entities.Activity, int64, error) {
	var activities []entities.Activity
	var count int64

	// Build query with filters
	query := repo.db.WithContext(ctx).Model(&entities.Activity{})
	
	// Apply filters
	query = repo.applyFilters(query, filter)

	// Count total records with filters
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and get results
	err := query.
		Preload("User").
		Preload("ActivityType").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Order("done_at DESC").
		Find(&activities).Error

	return activities, count, err
}

func (repo *activityRepo) applyFilters(query *gorm.DB, filter ActivityFilter) *gorm.DB {
	// Filter by user ID
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	// Filter by activity ID
	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}

	// Filter by activity type (join with activity_types table)
	if filter.ActivityType != nil {
		query = query.Joins("JOIN activity_types ON activities.activity_type_id = activity_types.id").
			Where("activity_types.name = ?", *filter.ActivityType)
	}

	// Filter by done_at date range
	if filter.DoneAtFrom != nil {
		query = query.Where("done_at >= ?", *filter.DoneAtFrom)
	}

	if filter.DoneAtTo != nil {
		query = query.Where("done_at <= ?", *filter.DoneAtTo)
	}

	// Filter by calories burned range
	if filter.CaloriesBurnedMin != nil {
		query = query.Where("calories_burned >= ?", *filter.CaloriesBurnedMin)
	}

	if filter.CaloriesBurnedMax != nil {
		query = query.Where("calories_burned <= ?", *filter.CaloriesBurnedMax)
	}

	return query
}

func (repo *activityRepo) GetActivity(ctx context.Context, filter ActivityFilter) (*entities.Activity, error) {
	var activity entities.Activity

	query := repo.db.WithContext(ctx).Model(&entities.Activity{})
	query = repo.applyFilters(query, filter)

	err := query.
		Preload("User").
		Preload("ActivityType").
		First(&activity).Error

	if err != nil {
		return nil, err
	}

	return &activity, nil
}

func (repo *activityRepo) CreateActivity(ctx context.Context, activity entities.Activity) (*entities.Activity, error) {
	err := repo.db.WithContext(ctx).Create(&activity).Error
	if err != nil {
		return nil, err
	}

	// Load relationships
	err = repo.db.WithContext(ctx).
		Preload("User").
		Preload("ActivityType").
		First(&activity, activity.ID).Error

	return &activity, err
}

func (repo *activityRepo) UpdateActivity(ctx context.Context, activity entities.Activity) (*entities.Activity, error) {
	err := repo.db.WithContext(ctx).Save(&activity).Error
	if err != nil {
		return nil, err
	}

	// Load relationships
	err = repo.db.WithContext(ctx).
		Preload("User").
		Preload("ActivityType").
		First(&activity, activity.ID).Error

	return &activity, err
}

func (repo *activityRepo) DeleteActivity(ctx context.Context, id uint) error {
	return repo.db.WithContext(ctx).Delete(&entities.Activity{}, id).Error
}

func (repo *activityRepo) GetActivityType(ctx context.Context, name string) (*entities.ActivityType, error) {
	var activityType entities.ActivityType

	err := repo.db.WithContext(ctx).
		Where("name = ?", name).
		First(&activityType).Error

	if err != nil {
		return nil, err
	}

	return &activityType, nil
}

func (repo *activityRepo) GetActivityTypes(ctx context.Context) ([]entities.ActivityType, error) {
	var activityTypes []entities.ActivityType

	err := repo.db.WithContext(ctx).
		Order("name").
		Find(&activityTypes).Error

	return activityTypes, err
}

func (repo *activityRepo) CreateActivityType(ctx context.Context, activityType entities.ActivityType) (*entities.ActivityType, error) {
	err := repo.db.WithContext(ctx).Create(&activityType).Error
	return &activityType, err
}