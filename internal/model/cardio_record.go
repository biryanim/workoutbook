package model

import (
	"github.com/pkg/errors"
	"time"
)

type CardioRecord struct {
	ID                int64
	WorkoutExerciseID int64
	DistanceKm        *float64
	DurationSeconds   int
	CreatedAt         time.Time
}

func NewCardioRecord(workoutExerciseID int64, durationsSeconds int, distanceKm *float64) *CardioRecord {
	return &CardioRecord{
		WorkoutExerciseID: workoutExerciseID,
		DistanceKm:        distanceKm,
		DurationSeconds:   durationsSeconds,
	}
}

func (c *CardioRecord) Validate() error {
	if c.DurationSeconds < 1 {
		return errors.New("duration must be positive")
	}
	if c.DistanceKm != nil && *c.DistanceKm < 0 {
		return errors.New("distance cannot be negative")
	}

	return nil
}

func (c *CardioRecord) CalculatePace() *float64 {
	if c.DistanceKm == nil || *c.DistanceKm == 0 {
		return nil
	}
	pace := float64(c.DurationSeconds) / 60.0 / *c.DistanceKm
	return &pace
}
