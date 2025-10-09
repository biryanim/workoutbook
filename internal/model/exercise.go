package model

import (
	"github.com/pkg/errors"
	"time"
)

const (
	ExerciseTypeStrength string = "strength"
	ExerciseTypeCardio   string = "cardio"
)

type Exercise struct {
	ID          int64
	Name        string
	Type        string
	MuscleGroup string
	Description string
	CreatedAt   time.Time
}

func (e *Exercise) IsStrength() bool {
	return e.Type == ExerciseTypeStrength
}

func (e *Exercise) IsCardio() bool {
	return e.Type == ExerciseTypeCardio
}

func (e *Exercise) Validate() error {
	if e.Name == "" {
		return errors.New("exercise name cannot be empty")
	}
	if e.Type != ExerciseTypeStrength && e.Type != ExerciseTypeCardio {
		return errors.New("invalid exercise type")
	}
	return nil
}
