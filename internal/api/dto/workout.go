package dto

import "time"

type Workout struct {
	ID     int64     `json:"id"`
	UserId int64     `json:"-"`
	Date   time.Time `json:"date"`
	Note   string    `json:"notes"`
	Name   string    `json:"name"`
}

type WorkoutExercise struct {
	WorkoutID  int64    `json:"workout_id"`
	ExerciseID int64    `json:"exercise_id"`
	Sets       int      `json:"sets"`
	Reps       int      `json:"reps,omitempty"`
	Weight     float64  `json:"weight,omitempty"`
	Duration   int      `json:"duration,omitempty"`
	Distance   float64  `json:"distance,omitempty"`
	Notes      string   `json:"notes"`
	Exercise   Exercise `json:"exercise"`
}

type Exercise struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	RecordType  string `json:"record_type"`
	MuscleGroup string `json:"muscle_group"`
	Description string `json:"description"`
}

type WorkoutExercises struct {
	Workout   Workout            `json:"workout"`
	Exercises []*WorkoutExercise `json:"exercises"`
}

type Pagination struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Limit     string `json:"limit"`
	Page      string `json:"page"`
}

type Record struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	ExerciseID int64     `json:"exercise_id"`
	Weight     float64   `json:"weight"`
	Reps       int       `json:"reps"`
	Date       time.Time `json:"date"`
	Notes      string    `json:"notes"`
	Duration   int       `json:"duration"`
	Sets       int       `json:"sets"`
	Distance   float64   `json:"distance"`
	Exercise   Exercise  `json:"exercise"`
}
