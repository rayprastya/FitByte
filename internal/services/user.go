package services

import (
	"context"
	"errors"
	"fitbyte/internal/entities"
	"fitbyte/internal/repositories"
	"fmt"
	"time"
)

type (
	UserService interface {
		GetUsersWithFilters(ctx context.Context, filter UserServiceFilter) (*entities.APIResponse, error)
		GetUser(ctx context.Context, id uint) (*entities.UserResponse, error)
		CreateUser(ctx context.Context, request entities.CreateUserRequest) (*entities.UserResponse, error)
		UpdateUser(ctx context.Context, id uint, request entities.UpdateUserRequest) (*entities.UserResponse, error)
		DeleteUser(ctx context.Context, id uint) error
	}

	userService struct {
		userRepository repositories.UserRepository
	}

	UserServiceParam struct {
		UserRepository repositories.UserRepository
	}

	UserServiceFilter struct {
		Limit    int
		Offset   int
		IsActive *bool
	}
)

func NewUserService(param UserServiceParam) UserService {
	return &userService{
		userRepository: param.UserRepository,
	}
}

func (svc *userService) GetUsersWithFilters(ctx context.Context, serviceFilter UserServiceFilter) (*entities.APIResponse, error) {
	repoFilter := repositories.UserFilter{
		Limit:    serviceFilter.Limit,
		Offset:   serviceFilter.Offset,
		IsActive: serviceFilter.IsActive,
	}

	users, total, err := svc.userRepository.GetUsersWithFilters(ctx, repoFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	userResponses := make([]entities.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = svc.mapToUserResponse(user)
	}

	return &entities.APIResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data: map[string]any{
			"users": userResponses,
			"meta": map[string]any{
				"total":  total,
				"limit":  serviceFilter.Limit,
				"offset": serviceFilter.Offset,
			},
		},
	}, nil
}

func (svc *userService) GetUser(ctx context.Context, id uint) (*entities.UserResponse, error) {
	filter := repositories.UserFilter{
		ID: &id,
	}

	user, err := svc.userRepository.GetUser(ctx, filter)
	if err != nil {
		return nil, errors.New("user not found")
	}

	response := svc.mapToUserResponse(*user)
	return &response, nil
}

func (svc *userService) CreateUser(ctx context.Context, request entities.CreateUserRequest) (*entities.UserResponse, error) {
	user := entities.User{
		Email:      request.Email,
		Name:       request.Name,
		Preference: request.Preference,
		WeightUnit: request.WeightUnit,
		HeightUnit: request.HeightUnit,
		Weight:     request.Weight,
		Height:     request.Height,
		ImageURI:   request.ImageURI,
	}

	createdUser, err := svc.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	response := svc.mapToUserResponse(*createdUser)
	return &response, nil
}

func (svc *userService) UpdateUser(ctx context.Context, id uint, request entities.UpdateUserRequest) (*entities.UserResponse, error) {
	filter := repositories.UserFilter{
		ID: &id,
	}

	existingUser, err := svc.userRepository.GetUser(ctx, filter)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if request.Email != nil {
		existingUser.Email = *request.Email
	}
	if request.Name != nil {
		existingUser.Name = request.Name
	}
	if request.Preference != nil {
		existingUser.Preference = request.Preference
	}
	if request.WeightUnit != nil {
		existingUser.WeightUnit = request.WeightUnit
	}
	if request.HeightUnit != nil {
		existingUser.HeightUnit = request.HeightUnit
	}
	if request.Weight != nil {
		existingUser.Weight = request.Weight
	}
	if request.Height != nil {
		existingUser.Height = request.Height
	}
	if request.ImageURI != nil {
		existingUser.ImageURI = request.ImageURI
	}

	updatedUser, err := svc.userRepository.UpdateUser(ctx, *existingUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := svc.mapToUserResponse(*updatedUser)
	return &response, nil
}

func (svc *userService) DeleteUser(ctx context.Context, id uint) error {
	return svc.userRepository.DeleteUser(ctx, id)
}

func (svc *userService) mapToUserResponse(user entities.User) entities.UserResponse {
	return entities.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		Preference: user.Preference,
		WeightUnit: user.WeightUnit,
		HeightUnit: user.HeightUnit,
		Weight:     user.Weight,
		Height:     user.Height,
		ImageURI:   user.ImageURI,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  user.UpdatedAt.Format(time.RFC3339),
	}
}