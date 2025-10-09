package main

import (
	"context"
	"fmt"
	authImpl "github.com/biryanim/workoutbook/internal/api/auth"
	"github.com/biryanim/workoutbook/internal/api/middleware"
	workoutImpl "github.com/biryanim/workoutbook/internal/api/workout"
	"github.com/biryanim/workoutbook/internal/client/db/pg"
	"github.com/biryanim/workoutbook/internal/client/db/transaction"
	"github.com/biryanim/workoutbook/internal/config"
	"github.com/biryanim/workoutbook/internal/config/env"
	userRepo "github.com/biryanim/workoutbook/internal/repository/user"
	workoutRepo "github.com/biryanim/workoutbook/internal/repository/workout"
	"github.com/biryanim/workoutbook/internal/service/auth"
	"github.com/biryanim/workoutbook/internal/service/workout"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	err := config.Load("local.env")
	if err != nil {
		log.Fatal(err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to load pg config: %v", err)
	}

	httpConfig, err := env.NewHTTPConfig()
	if err != nil {
		log.Fatalf("failed to load http config: %v", err)
	}

	jwtConfig, err := env.NewJWTConfig()
	if err != nil {
		log.Fatalf("failed to load jwt config: %v", err)
	}
	fmt.Println(pgConfig.DSN())
	dbClient, err := pg.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to initialize db client: %v", err)
	}
	err = dbClient.DB().Ping(ctx)
	if err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	txManager := transaction.NewTransactionManager(dbClient.DB())
	userRepository := userRepo.NewRepository(dbClient)
	workoutRepository := workoutRepo.NewRepository(dbClient)
	authService := auth.NewService(userRepository, txManager, jwtConfig)
	workoutService := workout.New(workoutRepository, txManager)
	authImpl := authImpl.NewImplementation(authService)
	workoutImpl := workoutImpl.NewImplementation(workoutService)

	r := gin.Default()

	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS())
	r.Use(middleware.XSSMiddleware())

	store := cookie.NewStore([]byte("your-secret-key-32-bytes-long!!"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24, // 24 часа
		HttpOnly: true,
		Secure:   false, // В проде сменить на true с HTTPS
		SameSite: http.SameSiteStrictMode,
	})
	r.Use(sessions.Sessions("workout_session", store))

	r.Use(csrf.Middleware(csrf.Options{
		Secret: "csrf-secret-key-change-in-production",
		ErrorFunc: func(c *gin.Context) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "CSRF token mismatch",
			})
			c.Abort()
		},
	}))
	setupRoutes(r, authImpl, workoutImpl)

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	if err = r.Run(httpConfig.Address()); err != nil {
		log.Fatal(err)
	}

}

func setupRoutes(r *gin.Engine, authHandler *authImpl.Implementation, workoutHandler *workoutImpl.Implementation) {
	// Отдаем HTML страницы
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"csrf_token": csrf.GetToken(c),
		})
	})
	r.GET("/login.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"csrf_token": csrf.GetToken(c),
		})
	})
	r.GET("/register.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"csrf_token": csrf.GetToken(c),
		})
	})
	r.GET("/dashboard.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"csrf_token": csrf.GetToken(c),
		})
	})
	r.GET("/workout-detail.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "workout-detail.html", gin.H{
			"csrf_token": csrf.GetToken(c),
		})
	})

	// Публичные API роуты
	public := r.Group("/api")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.GET("/csrf-token", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"csrf_token": csrf.GetToken(c)})
		})
	}

	// Защищенные API роуты
	protected := r.Group("/api")
	protected.Use(authHandler.AuthMiddleware())
	{
		protected.GET("/exercises", workoutHandler.ListExercises)

		protected.POST("/workouts", workoutHandler.CreateWorkout)
		protected.GET("/workouts", workoutHandler.ListWorkouts)
		protected.GET("/workouts/:id", workoutHandler.GetWorkout)
		protected.PUT("/workouts/:id", workoutHandler.UpdateWorkout)
		protected.DELETE("/workouts/:id", workoutHandler.DeleteWorkout)

		protected.POST("/workouts/:id/exercises", workoutHandler.AddExerciseToWorkout)
		protected.POST("/workouts/:id/cardio", workoutHandler.AddCardioToWorkout)

		protected.DELETE("/sets/:set_id", workoutHandler.DeleteExerciseSet)
		protected.DELETE("/cardio/:cardio_id", workoutHandler.DeleteCardioRecord)

		protected.GET("/records", workoutHandler.GetPersonalRecords)
	}
}
