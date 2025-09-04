package services

import (
	"context"
	"errors"
	"fitbyte/internal/entities"
	"fitbyte/internal/repositories"
	"fmt"
	"strconv"
	"time"
)

type (
	ActivityService interface {
		GetActivitiesWithFilters(ctx context.Context, userID uint, filter ActivityServiceFilter) (*entities.APIResponse, error)
		CreateActivity(ctx context.Context, request entities.CreateActivityRequest, userID uint) (*entities.ActivityResponse, error)
		GetActivityTypes(ctx context.Context) ([]entities.ActivityType, error)
	}

	activityService struct {
		activityRepository repositories.ActivityRepository
	}

	ActivityServiceParam struct {
		ActivityRepository repositories.ActivityRepository
	}

	ActivityServiceFilter struct {
		ActivityType      *string
		Limit             int
		Offset            int
		DoneAtFrom        *string
		DoneAtTo          *string
		CaloriesBurnedMin *float64
		CaloriesBurnedMax *float64
	}
)

func NewActivityService(param ActivityServiceParam) ActivityService {
	return &activityService{
		activityRepository: param.ActivityRepository,
	}
}

func (svc *activityService) GetActivitiesWithFilters(ctx context.Context, userID uint, serviceFilter ActivityServiceFilter) (*entities.APIResponse, error) {
	repoFilter := repositories.ActivityFilter{
		UserID:            &userID,
		ActivityType:      serviceFilter.ActivityType,
		Limit:             serviceFilter.Limit,
		Offset:            serviceFilter.Offset,
		CaloriesBurnedMin: serviceFilter.CaloriesBurnedMin,
		CaloriesBurnedMax: serviceFilter.CaloriesBurnedMax,
	}

	if serviceFilter.DoneAtFrom != nil {
		if doneAtFrom, err := time.Parse(time.RFC3339, *serviceFilter.DoneAtFrom); err == nil {
			repoFilter.DoneAtFrom = &doneAtFrom
		}
	}

	if serviceFilter.DoneAtTo != nil {
		if doneAtTo, err := time.Parse(time.RFC3339, *serviceFilter.DoneAtTo); err == nil {
			repoFilter.DoneAtTo = &doneAtTo
		}
	}

	activities, total, err := svc.activityRepository.GetActivitiesWithFilters(ctx, repoFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}

	activityResponses := make([]entities.ActivityResponse, len(activities))
	for i, activity := range activities {
		activityResponses[i] = svc.mapToActivityResponse(activity)
	}

	return &entities.APIResponse{
		Success: true,
		Message: "Activities retrieved successfully",
		Data: map[string]any{
			"activities": activityResponses,
			"meta": map[string]any{
				"total":  total,
				"limit":  serviceFilter.Limit,
				"offset": serviceFilter.Offset,
			},
		},
	}, nil
}

func (svc *activityService) CreateActivity(ctx context.Context, request entities.CreateActivityRequest, userID uint) (*entities.ActivityResponse, error) {
	activityType, err := svc.activityRepository.GetActivityType(ctx, request.ActivityType)
	if err != nil {
		return nil, errors.New("invalid activity type")
	}

	doneAt, err := time.Parse(time.RFC3339, request.DoneAt)
	if err != nil {
		return nil, errors.New("invalid date format, expected ISO 8601")
	}

	if request.DurationInMinutes < 1 {
		return nil, errors.New("duration must be at least 1 minute")
	}

	caloriesBurned := float64(request.DurationInMinutes) * activityType.CaloriesPerMinute

	activity := entities.Activity{
		UserID:            userID,
		ActivityTypeID:    activityType.ID,
		DurationInMinutes: request.DurationInMinutes,
		CaloriesBurned:    caloriesBurned,
		DoneAt:            doneAt,
		ActivityType:      *activityType,
	}

	createdActivity, err := svc.activityRepository.CreateActivity(ctx, activity)
	if err != nil {
		return nil, fmt.Errorf("failed to create activity: %w", err)
	}

	response := svc.mapToActivityResponse(*createdActivity)
	return &response, nil
}

func (svc *activityService) GetActivityTypes(ctx context.Context) ([]entities.ActivityType, error) {
	return svc.activityRepository.GetActivityTypes(ctx)
}

func (svc *activityService) mapToActivityResponse(activity entities.Activity) entities.ActivityResponse {
	return entities.ActivityResponse{
		ActivityID:        strconv.Itoa(int(activity.ID)),
		ActivityType:      activity.ActivityType.Name,
		DoneAt:            activity.DoneAt.Format(time.RFC3339),
		DurationInMinutes: activity.DurationInMinutes,
		CaloriesBurned:    activity.CaloriesBurned,
		CreatedAt:         activity.CreatedAt.Format(time.RFC3339),
	}
}
