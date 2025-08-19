package workout

import (
	"context"
	"fmt"
	"github.com/biryanim/workoutbook/internal/client/db"
	"github.com/biryanim/workoutbook/internal/model"
	"github.com/biryanim/workoutbook/internal/repository"
	"github.com/biryanim/workoutbook/internal/service"
	"github.com/jackc/pgx/v5"
)

var _ service.WorkoutService = (*serv)(nil)

type serv struct {
	workoutRepository repository.WorkoutRepository
	txManager         db.TxManager
}

func New(workoutRepository repository.WorkoutRepository, txManager db.TxManager) *serv {
	return &serv{
		workoutRepository: workoutRepository,
		txManager:         txManager,
	}
}

func (s *serv) CreateWorkout(ctx context.Context, workout *model.Workout) (int64, error) {
	id, err := s.workoutRepository.CreateWorkout(ctx, workout)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *serv) GetWorkouts(ctx context.Context, userId int64, pagination *model.WorkoutsFilter) ([]*model.Workout, error) {
	workouts, err := s.workoutRepository.ListWorkouts(ctx, userId, pagination)
	if err != nil {
		return nil, err
	}
	return workouts, nil
}

func (s *serv) GetWorkout(ctx context.Context, userId, workoutId int64) (*model.WorkoutExercises, error) {

	var (
		workout = &model.WorkoutExercises{}
		err     error
	)

	err = s.txManager.ReadCommited(ctx, func(ctx context.Context) error {

		workout.Workout, err = s.workoutRepository.GetWorkoutByID(ctx, workoutId, userId)
		if err != nil {
			return err
		}

		workout.Exercises, err = s.workoutRepository.GetExercisesByWorkoutID(ctx, workoutId)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (s *serv) AddExerciseToWorkout(ctx context.Context, userId int64, we *model.WorkoutExercise) error {
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		has, err := s.workoutRepository.IsUserHaveWorkout(ctx, userId, we.WorkoutID)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("workout not found for user %d", userId)
		}

		_, err = s.workoutRepository.AddWorkoutExercise(ctx, we)
		if err != nil {
			return err
		}

		err = s.UpdatePersonalRecord(ctx, userId, we.ExerciseID, we.Weight, we.Reps)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *serv) GetExercises(ctx context.Context, exerciseType string) ([]*model.Exercise, error) {
	exrs, err := s.workoutRepository.GetExercises(ctx, exerciseType)
	if err != nil {
		return nil, err
	}

	return exrs, nil
}

func (s *serv) UpdatePersonalRecord(ctx context.Context, userID, exerciseID int64, weight float64, reps int) error {
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		record, err := s.workoutRepository.GetPersonalRecord(ctx, userID, exerciseID)

		newMax := weight * (1 + float64(reps)/30)
		currentMax := record.Weight * (1 + float64(record.Reps)/30)

		user := &model.UserRecord{
			UserID:     userID,
			ExerciseID: exerciseID,
			Weight:     weight,
			Reps:       reps,
		}
		if err == pgx.ErrNoRows {
			_, err = s.workoutRepository.AddRecord(ctx, user)
			if err != nil {
				return err
			}
			return nil
		} else if newMax > currentMax {
			err = s.workoutRepository.UpdatePersonalRecord(ctx, user)
			if err != nil {
				return err
			}
			return nil
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (s *serv) GetPersonalRecords(ctx context.Context, userId int64) ([]*model.UserRecord, error) {
	records, err := s.workoutRepository.ListRecords(ctx, userId)
	if err != nil {
		return nil, err
	}

	return records, nil
}
