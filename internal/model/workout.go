package model

import (
	"database/sql"
	"time"
)

type Workout struct {
	ID        int64
	UserID    int64
	Date      time.Time
	Notes     string
	Name      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type WorkoutSet struct {
	ID         int64
	WorkoutID  int64
	ExerciseID int64
}

type WorkoutsFilter struct {
	StartDate time.Time
	EndDate   time.Time
	Offset    uint64
	Limit     uint64
}

type WorkoutExercise struct {
	ID         int64
	WorkoutID  int64
	ExerciseID int64
	Sets       int
	Reps       int
	Weight     float64
	Duration   int
	Distance   float64
	Exercise   Exercise
}

type Exercise struct {
	ID          int64
	Name        string
	Type        string
	MuscleGroup string
	Description string
}

type WorkoutExercises struct {
	Workout   *Workout
	Exercises []*WorkoutExercise
}

type UserRecord struct {
	ID         int64
	UserID     int64
	ExerciseID int64
	Weight     float64
	Reps       int
	Date       time.Time
	Notes      string
	Exercise   Exercise
}
