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
	"github.com/microcosm-cc/bluemonday"
	"log"
	"net/http"
)

func XSSMiddleware() gin.HandlerFunc {
	p := bluemonday.UGCPolicy()

	return func(c *gin.Context) {
		for _, values := range c.Request.URL.Query() {
			for i, value := range values {
				values[i] = p.Sanitize(value)
			}
		}

		// Санитизация form данных
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if err := c.Request.ParseForm(); err == nil {
				for _, values := range c.Request.PostForm {
					for i, value := range values {
						values[i] = p.Sanitize(value)
					}
				}
			}
		}

		c.Next()
	}
}

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

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
	//r.Use(SecurityHeadersMiddleware())
	//r.Use(XSSMiddleware())
	//
	//store := cookie.NewStore([]byte("your-secret-key-32-bytes-long!!"))
	//store.Options(sessions.Options{
	//	Path:     "/",
	//	MaxAge:   3600 * 24,               // 24 часа
	//	HttpOnly: true,                    // Защита от XSS
	//	Secure:   false,                   // Установите true в production с HTTPS
	//	SameSite: http.SameSiteStrictMode, // Дополнительная защита от CSRF
	//})
	//r.Use(sessions.Sessions("workout_session", store))
	//
	//r.Use(csrf.Middleware(csrf.Options{
	//	Secret: "csrf-secret-key-change-in-production",
	//	ErrorFunc: func(c *gin.Context) {
	//		c.JSON(http.StatusForbidden, gin.H{
	//			"error": "CSRF token mismatch",
	//		})
	//		c.Abort()
	//	},
	//}))
	//
	//r.GET("/api/csrf-token", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"csrf_token": csrf.GetToken(c),
	//	})
	//})

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

//r.Static("/static", "./static")
//r.LoadHTMLGlob("templates/*")
//r.GET("/", func(c *gin.Context) {
//	c.HTML(http.StatusOK, "index.html", nil)
//})
