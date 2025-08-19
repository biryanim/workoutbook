package main

import (
	"context"
	"fmt"
	authImpl "github.com/biryanim/workoutbook/internal/api/auth"
	workoutImpl "github.com/biryanim/workoutbook/internal/api/workout"
	"github.com/biryanim/workoutbook/internal/client/db/pg"
	"github.com/biryanim/workoutbook/internal/client/db/transaction"
	"github.com/biryanim/workoutbook/internal/config"
	"github.com/biryanim/workoutbook/internal/config/env"
	userRepo "github.com/biryanim/workoutbook/internal/repository/user"
	workoutRepo "github.com/biryanim/workoutbook/internal/repository/workout"
	"github.com/biryanim/workoutbook/internal/service/auth"
	"github.com/biryanim/workoutbook/internal/service/workout"
	"github.com/gin-gonic/gin"
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
	public := r.Group("/api")
	{
		public.POST("/register", authImpl.Register)
		public.POST("/login", authImpl.Login)
	}
	protected := r.Group("/api")
	protected.Use(authImpl.AuthMiddleware())
	{
		protected.GET("/exercises", workoutImpl.ListExercises)

		protected.POST("/workouts", workoutImpl.CreateWorkout)
		protected.GET("/workouts", workoutImpl.ListWorkouts)
		protected.GET("/workouts/:id", workoutImpl.GetWorkout)
		protected.POST("/workouts/:id/exercises", workoutImpl.AddExerciseToWorkout)

		protected.GET("/records", workoutImpl.GetPersonalRecords)
	}

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	if err = r.Run(httpConfig.Address()); err != nil {
		log.Fatal(err)
	}

}
