package model

import (
	"database/sql"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CreateUserParams struct {
	Email string
	Name  string
	//Age      int
	Password string
}

type LoginUserParams struct {
	Email    string
	Password string
}

type User struct {
	ID   int64
	Name string
	//Age       int
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserClaims struct {
	jwt.StandardClaims
	UserID int64
}
