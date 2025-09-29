package service

import (
	"context"
	"github.com/biryanim/workoutbook/internal/model"
	"time"
)

type AuthService interface {
	Register(ctx context.Context, userParams *model.CreateUserParams) (int64, error)
	Login(ctx context.Context, usrLoginParams *model.LoginUserParams) (*model.UserLoginResp, error)
	Check(ctx context.Context, token string) (int64, bool, error)
}

type WorkoutService interface {
	CreateWorkout(ctx context.Context, workout *model.Workout) (int64, error)
	GetWorkouts(ctx context.Context, userId int64, pagination *model.WorkoutsFilter) ([]*model.Workout, error)
	GetWorkout(ctx context.Context, userId, workoutId int64) (*model.WorkoutExercises, error)

	AddExerciseToWorkout(ctx context.Context, userId int64, we *model.WorkoutExercise) error
	GetExercises(ctx context.Context, exerciseType string) ([]*model.Exercise, error)

	UpdatePersonalRecord(ctx context.Context, userID int64, we *model.WorkoutExercise, date time.Time) error
	GetPersonalRecords(ctx context.Context, userId int64) ([]*model.UserRecord, error)
}
