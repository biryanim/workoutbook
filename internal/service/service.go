package service

import (
	"context"
	"github.com/biryanim/workoutbook/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, userParams *model.CreateUserParams) (int64, error)
	Login(ctx context.Context, usrLoginParams *model.LoginUserParams) (string, error)
	Check(ctx context.Context, token string) (int64, bool, error)
}

type WorkoutService interface {
	CreateWorkout(ctx context.Context, workout *model.Workout) (int64, error)
	GetWorkouts(ctx context.Context, userId int64, pagination *model.WorkoutsFilter) ([]*model.Workout, error)
	GetWorkout(ctx context.Context, userId, workoutId int64) (*model.WorkoutExercises, error)
	AddExerciseToWorkout(ctx context.Context, userId int64, we *model.WorkoutExercise) error
}
