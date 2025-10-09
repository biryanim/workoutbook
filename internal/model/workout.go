package model

import (
	"database/sql"
	"github.com/pkg/errors"
	"time"
)

type Workout struct {
	ID        int64
	UserID    int64
	Name      string
	Notes     string
	Date      time.Time
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	Exercises []*WorkoutExercise
}

func NewWorkout(userID int64, name, notes string, date time.Time) *Workout {
	return &Workout{
		UserID:    userID,
		Name:      name,
		Notes:     notes,
		Date:      date,
		Exercises: make([]*WorkoutExercise, 0),
	}
}

func (w *Workout) Validate() error {
	if w.UserID == 0 {
		return errors.New("invalid user ID")
	}
	if w.Name == "" {
		return errors.New("workout name cannot be empty")
	}
	return nil
}

func (w *Workout) AddExercise(exercise *WorkoutExercise) {
	w.Exercises = append(w.Exercises, exercise)
}

type WorkoutsFilter struct {
	StartDate time.Time
	EndDate   time.Time
	Offset    uint64
	Limit     uint64
}

type WorkoutExercises struct {
	Workout   *Workout
	Exercises []*WorkoutExercise
}

type UpdateWorkout struct {
	Name  *string
	Notes *string
	Date  *time.Time
}
