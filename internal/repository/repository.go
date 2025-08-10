package repository

import (
	"context"
	"github.com/biryanim/workoutbook/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.CreateUserParams) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type WorkoutRepository interface {
	Create(ctx context.Context, workout *model.Workout) (int64, error)
}
