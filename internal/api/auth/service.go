package auth

import (
	"fmt"
	"github.com/biryanim/workoutbook/internal/api/dto"
	"github.com/biryanim/workoutbook/internal/converter"
	apperrors "github.com/biryanim/workoutbook/internal/errors"
	"github.com/biryanim/workoutbook/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	auth = "Bearer "
)

type Implementation struct {
	authService service.AuthService
}

func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{authService: authService}
}

func (i *Implementation) Register(c *gin.Context) {
	var registerReq dto.UserRegisterRequest

	if err := c.ShouldBindJSON(&registerReq); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	resp, err := i.authService.Register(c.Request.Context(), converter.FromUserRegistrationRequest(&registerReq))
	if err != nil {
		fmt.Println(err)
		appErr := apperrors.FromError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user id": resp,
	})
}

func (i *Implementation) Login(c *gin.Context) {
	var loginReq dto.UserLoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	token, err := i.authService.Login(c.Request.Context(), converter.FromUserLoginRequest(&loginReq))
	if err != nil {
		fmt.Println(err)
		appErr := apperrors.FromError(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (i *Implementation) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "token required",
			})
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, auth)
		userId, access, err := i.authService.Check(c.Request.Context(), token)
		if err != nil || !access {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "token invalid",
			})
			c.Abort()
			return
		}
		c.Set("userID", userId)
		c.Next()
	}
}
