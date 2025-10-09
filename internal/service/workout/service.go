package workout

import (
	"context"
	"database/sql"
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

func (s *serv) CreateWorkout(ctx context.Context, workout *model.Workout) (*model.Workout, error) {
	_, err := s.workoutRepository.CreateWorkout(ctx, workout)
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (s *serv) GetWorkouts(ctx context.Context, userId int64, pagination *model.WorkoutsFilter) ([]*model.Workout, error) {
	workouts, err := s.workoutRepository.ListWorkouts(ctx, userId, pagination)
	if err != nil {
		return nil, err
	}
	return workouts, nil
}

func (s *serv) GetWorkout(ctx context.Context, userId, workoutId int64) (*model.Workout, error) {
	workout, err := s.workoutRepository.GetWorkoutByID(ctx, workoutId, userId)
	if err != nil {
		return nil, err
	}

	exercises, err := s.workoutRepository.GetWorkoutExercises(ctx, workoutId)
	if err != nil {
		return nil, err
	}

	for _, ex := range exercises {
		if ex.Exercise.IsStrength() {
			sets, err := s.workoutRepository.GetExerciseSets(ctx, ex.ID)
			if err != nil {
				return nil, err
			}
			ex.Sets = sets
		} else if ex.Exercise.IsCardio() {
			cardio, err := s.workoutRepository.GetCardioRecord(ctx, ex.ID)
			if err != nil {
				return nil, err
			}
			ex.Cardio = cardio
		}

	}

	workout.Exercises = exercises
	return workout, nil
}

func (s *serv) ListWorkouts(ctx context.Context, userID int64) ([]*model.Workout, error) {
	return s.workoutRepository.ListWorkoutsByUserID(ctx, userID)
}

func (s *serv) UpdateWorkout(ctx context.Context, workoutID, userId int64, update *model.UpdateWorkout) error {
	workout, err := s.workoutRepository.GetWorkoutByID(ctx, workoutID, userId)
	if err != nil {
		return err
	}

	fmt.Println("SERVICE: ", update.Date)
	if update.Name != nil {
		workout.Name = *update.Name
	}
	if update.Notes != nil {
		workout.Notes = *update.Notes
	}
	if update.Date != nil {
		workout.Date = *update.Date
	}
	workout.UpdatedAt = sql.NullTime{
		Time: time.Now(),
	}

	return s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		return s.workoutRepository.UpdateWorkout(ctx, workout)
	})
}

func (s *serv) DeleteWorkout(ctx context.Context, workoutID, userId int64) error {
	return s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		workout, err := s.workoutRepository.GetWorkoutByID(ctx, workoutID, userId)
		if err != nil {
			return err
		}
		if workout.UserID != userId {
			return errors.New("workout not found")
		}
		return s.workoutRepository.DeleteWorkout(ctx, workoutID)
	})
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
		_, err := s.workoutRepository.GetWorkoutByID(ctx, we.WorkoutID, userId)
		if err != nil {
			return err
		}

		if err = s.workoutRepository.AddWorkoutExercise(ctx, we); err != nil {
			return err
		}

		var (
			maxOneRM float64
			maxReps  int
		)

		for _, set := range we.Sets {
			set.WorkoutExerciseID = we.ID
			if err = set.Validate(); err != nil {
				return err
			}
			if err = s.workoutRepository.AddExerciseSet(ctx, set); err != nil {
				return err
			}

			oneRM := set.Calculate1RM()
			if oneRM > maxOneRM {
				maxOneRM = oneRM
			}

			if set.Reps > maxReps {
				maxReps = set.Reps
			}

		}

		if err := s.updatePersonalRecordIfBetter(ctx, userId, we.ExerciseID, model.RecordTypeMaxWeight, maxOneRM, &we.ID); err != nil {
			return err
		}
		if err := s.updatePersonalRecordIfBetter(ctx, userId, we.ExerciseID, model.RecordTypeMaxReps, float64(maxReps), &we.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *serv) DeleteExerciseSet(ctx context.Context, exerciseSetID int64) error {
	return s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		set, err := s.workoutRepository.GetExerciseSetByID(ctx, exerciseSetID)
		if err != nil {
			return err
		}

		workoutExerciseID := set.WorkoutExerciseID

		if err := s.workoutRepository.DeleteExerciseSet(ctx, exerciseSetID); err != nil {
			return err
		}
		sets, err := s.workoutRepository.GetExerciseSets(ctx, workoutExerciseID)
		if err != nil {
			return fmt.Errorf("failed to get remaining sets: %w", err)
		}

		if len(sets) == 0 {
			fmt.Println("AAAAAAAAAAAAAAAAAAAAAAa")
			err = s.workoutRepository.DeleteWorkoutExercise(ctx, workoutExerciseID)
			if err != nil {
				return fmt.Errorf("failed to delete workout exercise: %w", err)
			}
		} else {
			return s.workoutRepository.ReorderExerciseSets(ctx, workoutExerciseID)
		}

		return nil
	})
}

func (s *serv) AddCardioToWorkout(ctx context.Context, workoutID, userID, exerciseID int64, notes string, cardio *model.CardioRecord) error {
	_, err := s.workoutRepository.GetWorkoutByID(ctx, workoutID, userID)
	if err != nil {
		return err
	}

	exercise, err := s.workoutRepository.GetExerciseByID(ctx, exerciseID)
	if err != nil {
		return err
	}

	if exercise.Type != model.ExerciseTypeCardio {
		return errors.New("not cardio exercise")
	}

	workoutExercise := model.NewWorkoutExercise(workoutID, exerciseID, notes)
	err = s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		if err := s.workoutRepository.AddWorkoutExercise(ctx, workoutExercise); err != nil {
			return err
		}

		cardio.WorkoutExerciseID = workoutExercise.ID
		if err := cardio.Validate(); err != nil {
			return err
		}

		if err := s.workoutRepository.AddCardioRecord(ctx, cardio); err != nil {
			return err
		}

		if cardio.DistanceKm != nil && *cardio.DistanceKm > 0 {
			if err := s.updatePersonalRecordIfBetter(ctx, userID, exerciseID, model.RecordTypeMaxDistance, *cardio.DistanceKm, &workoutExercise.ID); err != nil {
				return err
			}
		}

		if cardio.DurationSeconds > 0 {
			if err := s.updatePersonalRecordIfBetter(ctx, userID, exerciseID, model.RecordTypeBestTime, float64(cardio.DurationSeconds), &workoutExercise.ID); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (s *serv) updatePersonalRecordIfBetter(ctx context.Context, userID, exerciseID int64, recordType model.RecordType, newValue float64, workoutExerciseID *int64) error {
	if newValue <= 0 {
		return nil
	}

	existingRecord, err := s.workoutRepository.GetPersonalRecord(ctx, userID, exerciseID, recordType)

	if err != nil || existingRecord == nil {
		record := model.NewPersonalRecord(userID, exerciseID, recordType, newValue)
		record.WorkoutExerciseID = workoutExerciseID
		return s.workoutRepository.UpsertPersonalRecord(ctx, record)
	}

	if existingRecord.IsNewRecord(newValue) {
		existingRecord.Value = newValue
		existingRecord.WorkoutExerciseID = workoutExerciseID
		existingRecord.AchievedAt = time.Now()
		return s.workoutRepository.UpsertPersonalRecord(ctx, existingRecord)
	}

	return nil
}

func (s *serv) GetPersonalRecords(ctx context.Context, userId int64) ([]*model.PersonalRecord, error) {
	return s.workoutRepository.GetPersonalRecords(ctx, userId)
}

func (s *serv) DeleteCardioRecord(ctx context.Context, cardioID int64) error {
	return s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		return s.workoutRepository.DeleteCardioRecord(ctx, cardioID)
	})
}

func (s *serv) AddSetToExercise(ctx context.Context, workoutExerciseID, userID int64, weight float64, reps int) error {
	return s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		currentSets, err := s.workoutRepository.GetExerciseSets(ctx, workoutExerciseID)
		if err != nil {
			return fmt.Errorf("failed to get current sets: %w", err)
		}

		newSetNumber := len(currentSets) + 1

		err = s.workoutRepository.CreateExerciseSet(ctx, workoutExerciseID, newSetNumber, weight, reps)
		if err != nil {
			return fmt.Errorf("failed to create set: %w", err)
		}

		return nil
	})
}
