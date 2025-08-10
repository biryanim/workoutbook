package auth

import (
	"context"
	"fmt"
	"github.com/pkg/errors"

	apperrors "github.com/biryanim/workoutbook/internal/errors"
	"github.com/biryanim/workoutbook/internal/model"

	"golang.org/x/crypto/bcrypt"
)

func (s *serv) Register(ctx context.Context, userParams *model.CreateUserParams) (int64, error) {
	existingUser, err := s.userRepository.GetByEmail(ctx, userParams.Email)
	if err != nil && !errors.Is(err, apperrors.ErrUserNotFound) {
		return 0, fmt.Errorf("failed to check existing auth: %w", err)
	}
	if existingUser != nil {
		return 0, apperrors.ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userParams.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	userParams.Password = string(passwordHash)

	id, err := s.userRepository.Create(ctx, userParams)
	if err != nil {
		return 0, fmt.Errorf("failed to create auth: %w", err)
	}

	return id, nil
}
