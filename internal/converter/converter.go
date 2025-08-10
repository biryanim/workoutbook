package converter

import (
	"github.com/biryanim/workoutbook/internal/api/dto"
	"github.com/biryanim/workoutbook/internal/model"
)

func FromUserRegistrationRequest(u *dto.UserRegisterRequest) *model.CreateUserParams {
	return &model.CreateUserParams{
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Username,
	}
}

func FromUserLoginRequest(u *dto.UserLoginRequest) *model.LoginUserParams {
	return &model.LoginUserParams{
		Email:    u.Email,
		Password: u.Password,
	}
}
