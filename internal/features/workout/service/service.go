package service

import (
	"context"

	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type WorkoutsService struct {
	workoutsRepository WorkoutsRepository
}

type WorkoutsRepository interface {
	GetWorkouts(
		ctx context.Context,
	//  
	) ([]domain.Workout, error)
	
	GetWorkout(
		ctx context.Context,
	// 
	) (domain.Workout, error)
}

func NewWorkoutsService(repo WorkoutsRepository) *WorkoutsService {
	return &WorkoutsService{
		workoutsRepository: repo,
	}
}