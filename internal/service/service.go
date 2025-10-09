package service

import (
	"context"
	"github.com/biryanim/workoutbook/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, userParams *model.CreateUserParams) (int64, error)
	Login(ctx context.Context, usrLoginParams *model.LoginUserParams) (*model.UserLoginResp, error)
	Check(ctx context.Context, token string) (int64, bool, error)
}

type WorkoutService interface {
	CreateWorkout(ctx context.Context, workout *model.Workout) (*model.Workout, error)
	GetWorkouts(ctx context.Context, userId int64, pagination *model.WorkoutsFilter) ([]*model.Workout, error)
	GetWorkout(ctx context.Context, userId, workoutId int64) (*model.Workout, error)
	ListWorkouts(ctx context.Context, userID int64) ([]*model.Workout, error)
	UpdateWorkout(ctx context.Context, workoutID, userId int64, update *model.UpdateWorkout) error
	DeleteWorkout(ctx context.Context, workoutID, userId int64) error

	AddExerciseToWorkout(ctx context.Context, userId int64, we *model.WorkoutExercise) error
	GetExercises(ctx context.Context, exerciseType string) ([]*model.Exercise, error)
	DeleteExerciseSet(ctx context.Context, exerciseSetID int64) error

	AddCardioToWorkout(ctx context.Context, workoutID, userID, exerciseID int64, notes string, cardio *model.CardioRecord) error
	//UpdatePersonalRecord(ctx context.Context, userID int64, we *model.WorkoutExercise, date time.Time) error
	GetPersonalRecords(ctx context.Context, userId int64) ([]*model.PersonalRecord, error)
	DeleteCardioRecord(ctx context.Context, cardioID int64) error
}
