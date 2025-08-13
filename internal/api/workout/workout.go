package workout

import (
	"fmt"
	"github.com/biryanim/workoutbook/internal/api/dto"
	"github.com/biryanim/workoutbook/internal/converter"
	apperrors "github.com/biryanim/workoutbook/internal/errors"
	"github.com/biryanim/workoutbook/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Implementation struct {
	workoutService service.WorkoutService
}

func NewImplementation(workoutService service.WorkoutService) *Implementation {
	return &Implementation{workoutService: workoutService}
}

func (i *Implementation) CreateWorkout(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var workout dto.Workout
	if err := c.ShouldBindJSON(&workout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	workout.UserId = userID

	id, err := i.workoutService.CreateWorkout(c.Request.Context(), converter.FromCreateWorkoutRequest(&workout))
	if err != nil {
		fmt.Println(err)
		appErr := apperrors.FromError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"workout_id": id})
}

func (i *Implementation) GetWorkout(c *gin.Context) {
	userID := c.GetInt64("user_id")
	workoutID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	workout, err := i.workoutService.GetWorkout(c.Request.Context(), userID, workoutID)
	if err != nil {
		fmt.Println(err)
		appErr := apperrors.FromError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
		return
	}
	fmt.Println(workout.Workout)
	resp := converter.ToGetWorkoutResp(workout)
	fmt.Println(resp)
	c.JSON(http.StatusOK, resp)
}

func (i *Implementation) ListWorkouts(c *gin.Context) {
	userID := c.GetInt64("user_id")

	pagination := dto.Pagination{
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
		Limit:     c.Query("limit"),
		Page:      c.Query("page"),
	}

	filter, err := converter.FromPaginationToFilter(&pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	workouts, err := i.workoutService.GetWorkouts(c.Request.Context(), userID, filter)
	if err != nil {
		fmt.Println(err)
		appErr := apperrors.FromError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, converter.ToWorkoutsResp(workouts))
}

func (i *Implementation) AddExerciseToWorkout(c *gin.Context) {
	userID := c.GetInt64("user_id")
	workoutID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var exerc dto.WorkoutExercise
	if err := c.ShouldBindJSON(&exerc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	exerc.WorkoutID = workoutID

	err = i.workoutService.AddExerciseToWorkout(c.Request.Context(), userID, converter.FromAddExerciseToWorkout(&exerc))
	if err != nil {
		fmt.Println(err)
		appErr := apperrors.FromError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"workout_id": workoutID})
}

func (i *Implementation) ListExercises(c *gin.Context) {
	exerciseType := c.Query("type")

	exercises, err := i.workoutService.GetExercises(c.Request.Context(), exerciseType)
	if err != nil {
		fmt.Println(err)
		appErr := apperrors.FromError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, converter.ToListExercisesResp(exercises))
}
