package auth

import (
	"context"
	"fmt"
	apperrors "github.com/biryanim/workoutbook/internal/errors"
	"github.com/biryanim/workoutbook/internal/model"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *serv) Login(ctx context.Context, userLogin *model.LoginUserParams) (string, error) {
	user, err := s.userRepository.GetByEmail(ctx, userLogin.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return "", apperrors.ErrInvalidCredentials
		}
		return "", fmt.Errorf("failed to get auth: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		return "", apperrors.ErrInvalidCredentials
	}

	token, err := generateToken(user.ID, s.jwtConfig.TokenSecret(), s.jwtConfig.TokenExpiration())
	if err != nil {
		return "", errors.Wrap(err, "failed to generate token")
	}

	return token, nil
}
