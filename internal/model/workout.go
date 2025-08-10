package model

import "time"

type Workout struct {
	ID     int64
	UserID int64
	Date   time.Time
	Notes  string
	Name   string
}

type WorkoutSet struct {
	ID         int64
	WorkoutID  int
	ExerciseID int
}
