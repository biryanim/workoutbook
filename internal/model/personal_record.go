package model

import (
	"github.com/pkg/errors"
	"time"
)

type RecordType string

const (
	RecordTypeMaxWeight   RecordType = "max_weight"
	RecordTypeMaxReps     RecordType = "max_reps"
	RecordTypeMaxDistance RecordType = "max_distance"
	RecordTypeBestTime    RecordType = "best_time"
)

type PersonalRecord struct {
	ID                int64
	UserID            int64
	ExerciseID        int64
	RecordType        RecordType
	Value             float64
	WorkoutExerciseID *int64
	AchievedAt        time.Time
	Exercise          *Exercise
}

func NewPersonalRecord(userID, exerciseID int64, recordType RecordType, value float64) *PersonalRecord {
	return &PersonalRecord{
		UserID:     userID,
		ExerciseID: exerciseID,
		RecordType: RecordType(recordType),
		Value:      value,
		AchievedAt: time.Now(),
	}
}

func (pr *PersonalRecord) Validate() error {
	if pr.UserID == 0 {
		return errors.New("invalid user ID")
	}
	if pr.ExerciseID == 0 {
		return errors.New("invalid exercise ID")
	}
	if pr.Value <= 0 {
		return errors.New("record value must be positive")
	}
	return nil
}

func (pr *PersonalRecord) IsNewRecord(newValue float64) bool {
	switch pr.RecordType {
	case RecordTypeMaxWeight, RecordTypeMaxReps, RecordTypeMaxDistance:
		return newValue > pr.Value
	case RecordTypeBestTime:
		return newValue < pr.Value
	default:
		return false
	}
}
