package auth

import (
	"github.com/biryanim/workoutbook/internal/client/db"
	"github.com/biryanim/workoutbook/internal/config"
	"github.com/biryanim/workoutbook/internal/repository"
	"github.com/biryanim/workoutbook/internal/service"
)

var _ service.AuthService = (*serv)(nil)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
	jwtConfig      config.JWTConfig
}

func NewService(userRepository repository.UserRepository, txManager db.TxManager, jwtConfig config.JWTConfig) service.AuthService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
		jwtConfig:      jwtConfig,
	}
}
