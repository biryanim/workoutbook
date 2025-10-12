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
	// Workout
	CreateWorkout(ctx context.Context, workout *model.Workout) (int64, error)
	GetWorkoutByID(ctx context.Context, workoutID, userId int64) (*model.Workout, error)
	ListWorkoutsByUserID(ctx context.Context, userId int64) ([]*model.Workout, error)
	ListWorkouts(ctx context.Context, userId int64, filter *model.WorkoutsFilter) ([]*model.Workout, error)
	UpdateWorkout(ctx context.Context, workout *model.Workout) error
	DeleteWorkout(ctx context.Context, id int64) error

	AddWorkoutExercise(ctx context.Context, we *model.WorkoutExercise) error
	GetWorkoutExercises(ctx context.Context, workoutID int64) ([]*model.WorkoutExercise, error)
	DeleteWorkoutExercise(ctx context.Context, workoutExerciseID int64) error

	AddExerciseSet(ctx context.Context, set *model.ExerciseSet) error
	GetExerciseSets(ctx context.Context, workoutExerciseID int64) ([]*model.ExerciseSet, error)
	GetExerciseSetByID(ctx context.Context, id int64) (*model.ExerciseSet, error)
	UpdateExerciseSet(ctx context.Context, set *model.ExerciseSet) error
	DeleteExerciseSet(ctx context.Context, id int64) error
	ReorderExerciseSets(ctx context.Context, workoutExerciseID int64) error

	GetExercises(ctx context.Context, typ string) ([]*model.Exercise, error)
	GetExerciseByID(ctx context.Context, exerciseID int64) (*model.Exercise, error)
	GetPersonalRecord(ctx context.Context, userID, exerciseID int64, recordType model.RecordType) (*model.PersonalRecord, error)
	GetPersonalRecords(ctx context.Context, userID int64) ([]*model.PersonalRecord, error)
	UpsertPersonalRecord(ctx context.Context, record *model.PersonalRecord) error

	AddCardioRecord(ctx context.Context, cardio *model.CardioRecord) error
	GetCardioRecordByWorkoutExerciseID(ctx context.Context, workoutExerciseID int64) (*model.CardioRecord, error)
	GetCardioRecordByID(ctx context.Context, id int64) (*model.CardioRecord, error)
	UpdateCardioRecord(ctx context.Context, cardio *model.CardioRecord) error
	DeleteCardioRecord(ctx context.Context, id int64) error

	CreateExerciseSet(ctx context.Context, workoutExerciseID int64, setNumber int, weight float64, reps int) error
}
