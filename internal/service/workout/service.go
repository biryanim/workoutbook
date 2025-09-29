package workout

import (
	"context"
	"fmt"
	"github.com/biryanim/workoutbook/internal/client/db"
	"github.com/biryanim/workoutbook/internal/model"
	"github.com/biryanim/workoutbook/internal/repository"
	"github.com/biryanim/workoutbook/internal/service"
	"github.com/pkg/errors"
	"time"
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

func (s *serv) GetExercises(ctx context.Context, exerciseType string) ([]*model.Exercise, error) {
	exrs, err := s.workoutRepository.GetExercises(ctx, exerciseType)
	if err != nil {
		return nil, err
	}

	return exrs, nil
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

		date, err := s.workoutRepository.AddWorkoutExercise(ctx, we)
		if err != nil {
			return err
		}

		err = s.UpdatePersonalRecord(ctx, userId, we, date)
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

func (s *serv) UpdatePersonalRecord(ctx context.Context, userID int64, we *model.WorkoutExercise, date time.Time) error {
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		record, err := s.workoutRepository.GetPersonalRecord(ctx, userID, we.ExerciseID)
		user := &model.UserRecord{
			UserID:     userID,
			ExerciseID: we.ExerciseID,
			Weight:     we.Weight,
			Reps:       we.Reps,
			Sets:       we.Sets,
			Duration:   we.Duration,
			Distance:   we.Distance,
			Date:       date,
		}
		if record == nil {
			_, err = s.workoutRepository.AddRecord(ctx, user)
			if err != nil {
				return err
			}
			return nil
		}

		exercise, err := s.workoutRepository.GetExerciseByID(ctx, we.ExerciseID)
		if err != nil {
			return err
		}
		we.Exercise = *exercise

		switch we.Exercise.Type {
		case "strength":
			// Особая логика для упражнений с собственным весом без веса
			switch we.Exercise.Name {
			case "Отжимания", "Подтягивания", "Скручивания":
				if we.Reps > record.Reps {
					err = s.workoutRepository.UpdatePersonalRecord(ctx, user)
					if err != nil {
						return err
					}
				}
			default:
				newMax := we.Weight * (1 + float64(we.Reps)/30)
				currentMax := record.Weight * (1 + float64(record.Reps)/30)
				if newMax > currentMax {
					err = s.workoutRepository.UpdatePersonalRecord(ctx, user)
					if err != nil {
						return err
					}
				}
			}
		case "cardio":
			switch we.Exercise.RecordType {
			case "distance":
				if we.Distance > record.Distance {
					err = s.workoutRepository.UpdatePersonalRecord(ctx, user)
					if err != nil {
						return err
					}
				}
			case "duration":
				if we.Duration > record.Duration {
					err = s.workoutRepository.UpdatePersonalRecord(ctx, user)
					if err != nil {
						return err
					}
				}
			case "sets":
				if we.Sets > record.Sets {
					err = s.workoutRepository.UpdatePersonalRecord(ctx, user)
					if err != nil {
						return err
					}
				}
			default:
				return errors.New("unknown cardio record type")
			}
		default:
			return errors.New("unsupported exercise type")
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
