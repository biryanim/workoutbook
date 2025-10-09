package model

import (
	"github.com/pkg/errors"
	"time"
)

type ExerciseSet struct {
	ID                int64
	WorkoutExerciseID int64
	SetNumber         int
	Weight            float64
	Reps              int
	CreatedAt         time.Time
}

func NewExerciseSet(workoutExerciseID int64, setNumber int, weight float64, reps int) *ExerciseSet {
	return &ExerciseSet{
		WorkoutExerciseID: workoutExerciseID,
		SetNumber:         setNumber,
		Weight:            weight,
		Reps:              reps,
	}
}

func (s *ExerciseSet) Validate() error {
	if s.SetNumber < 1 {
		return errors.New("set number must be positive")
	}
	if s.Weight < 0 {
		return errors.New("weight cannot be negative")
	}
	if s.Reps < 1 {
		return errors.New("reps must be positive")
	}
	return nil
}

func (s *ExerciseSet) Calculate1RM() float64 {
	if s.Reps == 1 {
		return s.Weight
	}
	return s.Weight * (36.0 / (37.0 - float64(s.Reps)))
}

type UpdateExerciseSet struct {
	WorkoutExerciseID int64
	SetNumber         int
	Weight            float64
	Reps              int
	CreatedAt         time.Time
}
