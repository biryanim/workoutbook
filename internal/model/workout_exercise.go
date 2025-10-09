package model

import "time"

type WorkoutExercise struct {
	ID         int64
	WorkoutID  int64
	ExerciseID int64
	Notes      string
	Exercise   *Exercise
	CreatedAt  time.Time
	Sets       []*ExerciseSet
	Cardio     *CardioRecord
}

func NewWorkoutExercise(workoutID, exerciseID int64, notes string) *WorkoutExercise {
	return &WorkoutExercise{
		WorkoutID:  workoutID,
		ExerciseID: exerciseID,
		Notes:      notes,
		Sets:       make([]*ExerciseSet, 0),
	}
}

func (we *WorkoutExercise) CalculateVolume() float64 {
	var volume float64
	for _, set := range we.Sets {
		volume += set.Weight * float64(set.Reps)
	}
	return volume
}

func (we *WorkoutExercise) AddSet(set *ExerciseSet) {
	we.Sets = append(we.Sets, set)
}
