package model

import (
	"database/sql"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CreateUserParams struct {
	Email    string
	Password string
}

type LoginUserParams struct {
	Email    string
	Password string
}

type User struct {
	ID        int64
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserClaims struct {
	jwt.StandardClaims
	UserID int64
}

type UserLoginResp struct {
	Token string
}
