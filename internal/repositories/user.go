package repositories

import (
	"context"
	"fitbyte/internal/entities"

	"gorm.io/gorm"
)

type (
	UserRepository interface {
		GetUsersWithFilters(ctx context.Context, filter UserFilter) ([]entities.User, int64, error)
		GetUser(ctx context.Context, filter UserFilter) (*entities.User, error)
		CreateUser(ctx context.Context, user entities.User) (*entities.User, error)
		UpdateUser(ctx context.Context, user entities.User) (*entities.User, error)
		DeleteUser(ctx context.Context, id uint) error
	}

	userRepo struct {
		db *gorm.DB
	}

	UserRepoParam struct {
		DB *gorm.DB
	}

	UserFilter struct {
		ID       *uint
		Email    *string
		IsActive *bool
		Limit    int
		Offset   int
	}
)

func NewUserRepository(param UserRepoParam) UserRepository {
	return &userRepo{
		db: param.DB,
	}
}

func (repo *userRepo) GetUsersWithFilters(ctx context.Context, filter UserFilter) ([]entities.User, int64, error) {
	var users []entities.User
	var count int64

	query := repo.db.WithContext(ctx).Model(&entities.User{})
	
	query = repo.applyFilters(query, filter)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Limit(filter.Limit).
		Offset(filter.Offset).
		Order("created_at DESC").
		Find(&users).Error

	return users, count, err
}

func (repo *userRepo) applyFilters(query *gorm.DB, filter UserFilter) *gorm.DB {
	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}

	if filter.Email != nil {
		query = query.Where("email = ?", *filter.Email)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	return query
}

func (repo *userRepo) GetUser(ctx context.Context, filter UserFilter) (*entities.User, error) {
	var user entities.User

	query := repo.db.WithContext(ctx).Model(&entities.User{})
	query = repo.applyFilters(query, filter)

	err := query.First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepo) CreateUser(ctx context.Context, user entities.User) (*entities.User, error) {
	err := repo.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepo) UpdateUser(ctx context.Context, user entities.User) (*entities.User, error) {
	err := repo.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepo) DeleteUser(ctx context.Context, id uint) error {
	return repo.db.WithContext(ctx).Model(&entities.User{}).Where("id = ?", id).Update("is_active", false).Error
}