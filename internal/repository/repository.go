package repository

import (
	"context"
	"github.com/biryanim/workoutbook/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.CreateUserParams) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type WorkoutRepository interface {
	CreateWorkout(ctx context.Context, workout *model.Workout) (int64, error)
	GetWorkoutByID(ctx context.Context, workoutID, userId int64) (*model.Workout, error)
	ListWorkouts(ctx context.Context, userId int64, filter *model.WorkoutsFilter) ([]*model.Workout, error)
	AddWorkoutExercise(ctx context.Context, we *model.WorkoutExercise) (int64, error)
	GetExercisesByWorkoutID(ctx context.Context, workoutID int64) ([]*model.WorkoutExercise, error)
	IsUserHaveWorkout(ctx context.Context, userId, workoutId int64) (bool, error)
	GetExercises(ctx context.Context, typ string) ([]*model.Exercise, error)

	GetPersonalRecord(ctx context.Context, userID, exerciseID int64) (*model.UserRecord, error)
	AddRecord(ctx context.Context, user *model.UserRecord) (int64, error)
	UpdatePersonalRecord(ctx context.Context, user *model.UserRecord) error
	ListRecords(ctx context.Context, userId int64) ([]*model.UserRecord, error)
}
