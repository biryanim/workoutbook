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
		deletedSetNumber := set.SetNumber

		if err := s.workoutRepository.DeleteExerciseSet(ctx, exerciseSetID); err != nil {
			return err
		}

		return s.workoutRepository.ReorderExerciseSets(ctx, workoutExerciseID, deletedSetNumber)
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

//func (s *serv) UpdatePersonalRecord(ctx context.Context, userID int64, we *model.WorkoutExercise, date time.Time) error {
//	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
//		record, err := s.workoutRepository.GetPersonalRecord(ctx, userID, we.ExerciseID)
//		user := &model.UserRecord{
//			UserID:     userID,
//			ExerciseID: we.ExerciseID,
//			Weight:     we.Weight,
//			Reps:       we.Reps,
//			Sets:       we.Sets,
//			Duration:   we.Duration,
//			Distance:   we.Distance,
//			Date:       date,
//		}
//		if record == nil {
//			_, err = s.workoutRepository.AddRecord(ctx, user)
//			if err != nil {
//				return err
//			}
//			return nil
//		}
//
//		exercise, err := s.workoutRepository.GetExerciseByID(ctx, we.ExerciseID)
//		if err != nil {
//			return err
//		}
//		we.Exercise = *exercise
//
//		switch we.Exercise.Type {
//		case "strength":
//			// Особая логика для упражнений с собственным весом без веса
//			switch we.Exercise.Name {
//			case "Отжимания", "Подтягивания", "Скручивания":
//				if we.Reps > record.Reps {
//					err = s.workoutRepository.UpdatePersonalRecord(ctx, user)
//					if err != nil {
//						return err
//					}
//				}
//			default:
//				newMax := we.Weight * (1 + float64(we.Reps)/30)
//				currentMax := record.Weight * (1 + float64(record.Reps)/30)
//				if newMax > currentMax {
//					err = s.workoutRepository.UpdatePersonalRecord(ctx, user)
//					if err != nil {
//						return err
//					}
//				}
//			}
//		case "cardio":
//			switch we.Exercise.RecordType {
//			case "distance":
//				if we.Distance > record.Distance {
//					err = s.workoutRepository.UpdatePersonalRecord(ctx, user)
//					if err != nil {
//						return err
//					}
//				}
//			case "duration":
//				if we.Duration > record.Duration {
//					err = s.workoutRepository.UpdatePersonalRecord(ctx, user)
//					if err != nil {
//						return err
//					}
//				}
//			case "sets":
//				if we.Sets > record.Sets {
//					err = s.workoutRepository.UpdatePersonalRecord(ctx, user)
//					if err != nil {
//						return err
//					}
//				}
//			default:
//				return errors.New("unknown cardio record type")
//			}
//		default:
//			return errors.New("unsupported exercise type")
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		return err
//	}
//	return nil
//}
