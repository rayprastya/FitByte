package server

import (
	"fitbyte/internal/database"
	"fitbyte/internal/handlers"
	"fitbyte/internal/repositories"
	"fitbyte/internal/services"
)

type Handlers struct {
	HealthHandler   *handlers.HealthHandler
	UserHandler     *handlers.UserHandler
	ActivityHandler *handlers.ActivityHandler
}

func NewHandlers() *Handlers {
	healthHandler := handlers.NewHealthHandler()

	userRepo := repositories.NewUserRepository(repositories.UserRepoParam{
		DB: database.DB,
	})
	userService := services.NewUserService(services.UserServiceParam{
		UserRepository: userRepo,
	})
	userHandler := handlers.NewUserHandler(userService)

	activityRepo := repositories.NewActivityRepository(repositories.ActivityRepoParam{
		DB: database.DB,
	})
	activityService := services.NewActivityService(services.ActivityServiceParam{
		ActivityRepository: activityRepo,
	})
	activityHandler := handlers.NewActivityHandler(activityService)

	return &Handlers{
		HealthHandler:   healthHandler,
		UserHandler:     userHandler,
		ActivityHandler: activityHandler,
	}
}