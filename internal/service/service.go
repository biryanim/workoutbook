package service

import (
	"context"
	"github.com/biryanim/workoutbook/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, userParams *model.CreateUserParams) (int64, error)
	Login(ctx context.Context, usrLoginParams *model.LoginUserParams) (string, error)
	Check(ctx context.Context, token string) (bool, error)
}

type WorkoutService interface {
	CreateWorkout(ctx context.Context, )
}
