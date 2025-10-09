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
	var workout dto.CreateWorkoutRequest
	if err := c.ShouldBindJSON(&workout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	workout.UserID = userID
	w, err := converter.FromCreateWorkoutRequest(&workout)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	resp, err := i.workoutService.CreateWorkout(c.Request.Context(), w)
	if err != nil {
		fmt.Println(err)
		appErr := apperrors.FromError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"workout": converter.ToCreateWorkoutsResp(resp)})
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
	resp := converter.ToGetWorkoutResp(workout)
	c.JSON(http.StatusOK, gin.H{"workout": resp})
}

func (i *Implementation) ListWorkouts(c *gin.Context) {
	userID := c.GetInt64("user_id")

	workouts, err := i.workoutService.ListWorkouts(c.Request.Context(), userID)
	if err != nil {
		fmt.Println(err)
		appErr := apperrors.FromError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"workouts": converter.ToWorkoutsResp(workouts)})
}

func (i *Implementation) UpdateWorkout(c *gin.Context) {
	userID := c.GetInt64("user_id")
	workoutID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout id"})
		return
	}

	var req dto.UpdateWorkoutRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("API:", req.WorkoutDate)

	err = i.workoutService.UpdateWorkout(c.Request.Context(), workoutID, userID, converter.FromUpdateWorkoutRequest(&req))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "workout updated"})
}

func (i *Implementation) DeleteWorkout(c *gin.Context) {
	userID := c.GetInt64("user_id")
	workoutID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout id"})
		return
	}

	err = i.workoutService.DeleteWorkout(c.Request.Context(), workoutID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "workout deleted"})
}

func (i *Implementation) AddExerciseToWorkout(c *gin.Context) {
	userID := c.GetInt64("user_id")
	workoutID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.AddExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = i.workoutService.AddExerciseToWorkout(c.Request.Context(), userID, converter.FromAddExerciseToWorkout(&req, workoutID))
	if err != nil {
		fmt.Println(err)
		appErr := apperrors.FromError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "exercise added to workout"})
}

func (i *Implementation) AddCardioToWorkout(c *gin.Context) {
	userID := c.GetInt64("user_id")

	workoutID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout id"})
	}

	var req dto.AddCardioRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cardio := converter.FromAddCardioToWorkoutRequest(&req)

	err = i.workoutService.AddCardioToWorkout(c.Request.Context(), workoutID, userID, req.ExerciseID, req.Notes, cardio)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "cardio added to workout"})
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

	c.JSON(http.StatusOK, gin.H{"exercises": converter.ToListExercisesResp(exercises)})
}

func (i *Implementation) GetPersonalRecords(c *gin.Context) {
	userID := c.GetInt64("user_id")
	records, err := i.workoutService.GetPersonalRecords(c.Request.Context(), userID)
	if err != nil {
		fmt.Println(err)
		appErr := apperrors.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"records": converter.ToPersonalRecord(records)})
}

func (i *Implementation) DeleteExerciseSet(c *gin.Context) {
	//userID := c.GetInt64("user_id")
	setID, err := strconv.ParseInt(c.Param("set_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid set id"})
		return
	}

	err = i.workoutService.DeleteExerciseSet(c.Request.Context(), setID)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "exercise set deleted"})
}

func (i *Implementation) DeleteCardioRecord(c *gin.Context) {
	cardioID, err := strconv.ParseInt(c.Param("cardio_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cardio id"})
		return
	}

	err = i.workoutService.DeleteCardioRecord(c.Request.Context(), cardioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cardio id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cardio deleted"})
}
