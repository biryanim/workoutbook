package dto

import (
	"time"
)

type CreateWorkoutRequest struct {
	UserID      int64  `json:"-"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	WorkoutDate string `json:"workout_date"`
}

type Exercise struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
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

type Workout struct {
	ID        int64              `json:"id"`
	UserID    int64              `json:"user_id"`
	Name      string             `json:"name"`
	Notes     string             `json:"description"`
	Date      time.Time          `json:"workout_date"`
	Exercises []*WorkoutExercise `json:"exercises"`
}

type WorkoutExercise struct {
	ID         int64          `json:"id"`
	WorkoutID  int64          `json:"workout_id"`
	ExerciseID int64          `json:"exercise_id"`
	Notes      string         `json:"notes,omitempty"`
	Exercise   *Exercise      `json:"exercise,omitempty"`
	Sets       []*ExerciseSet `json:"sets,omitempty"`
	Cardio     *CardioRecord  `json:"cardio,omitempty"`
}

type ExerciseSet struct {
	ID                int64   `json:"id"`
	WorkoutExerciseID int64   `json:"workout_exercise_id"`
	SetNumber         int     `json:"set_number"`
	Weight            float64 `json:"weight"`
	Reps              int     `json:"reps"`
}

type CardioRecord struct {
	ID                int64    `json:"id"`
	WorkoutExerciseID int64    `json:"workout_exercise_id"`
	DistanceKm        *float64 `json:"distance_km"`
	DurationSeconds   int      `json:"duration_seconds"`
}

type UpdateWorkoutRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	WorkoutDate *string `json:"workout_date,omitempty"`
}

type AddExerciseRequest struct {
	ExerciseID int64    `json:"exercise_id" binding:"required"`
	Notes      string   `json:"notes"`
	Sets       []SetDTO `json:"sets"`
}

type SetDTO struct {
	SetNumber int     `json:"set_number" binding:"required"`
	Weight    float64 `json:"weight" binding:"required"`
	Reps      int     `json:"reps" binding:"required"`
}

type PersonalRecord struct {
	ID                int64     `json:"id"`
	UserID            int64     `json:"user_id"`
	ExerciseID        int64     `json:"exercise_id"`
	RecordType        string    `json:"record_type"`
	Value             float64   `json:"value"`
	WorkoutExerciseID *int64    `json:"workout_exercise_id"`
	AchievedAt        time.Time `json:"achieved_at"`
	Exercise          *Exercise `json:"exercise"`
}

type AddCardioRequest struct {
	ExerciseID      int64    `json:"exercise_id" binding:"required"`
	Notes           string   `json:"notes"`
	DistanceKm      *float64 `json:"distance_km"`
	DurationSeconds int      `json:"duration_seconds" binding:"required"`
}
