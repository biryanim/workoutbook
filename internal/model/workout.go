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
	WorkoutID  int
	ExerciseID int
}

type WorkoutsFilter struct {
	StartDate time.Time
	EndDate   time.Time
	Offset    uint64
	Limit     uint64
}
