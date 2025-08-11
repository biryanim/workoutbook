package converter

import (
	"github.com/biryanim/workoutbook/internal/api/dto"
	"github.com/biryanim/workoutbook/internal/model"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

func FromUserRegistrationRequest(u *dto.UserRegisterRequest) *model.CreateUserParams {
	return &model.CreateUserParams{
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Username,
	}
}

func FromUserLoginRequest(u *dto.UserLoginRequest) *model.LoginUserParams {
	return &model.LoginUserParams{
		Email:    u.Email,
		Password: u.Password,
	}
}

func FromCreateWorkoutRequest(r *dto.Workout) *model.Workout {
	return &model.Workout{
		UserID: r.UserId,
		Date:   r.Date,
		Name:   r.Name,
		Notes:  r.Note,
	}
}

func ToGetWorkoutResp(w *model.WorkoutExercises) *dto.WorkoutExercises {
	var (
		exercises []*dto.WorkoutExercise
	)
	for _, ex := range w.Exercises {
		dtoEx := &dto.WorkoutExercise{
			Sets:     ex.Sets,
			Reps:     ex.Reps,
			Weight:   ex.Weight,
			Duration: ex.Duration,
			Distance: ex.Distance,
			Exercise: dto.Exercise{
				Name:        ex.Exercise.Name,
				Type:        ex.Exercise.Type,
				MuscleGroup: ex.Exercise.MuscleGroup,
				Description: ex.Exercise.Description,
			},
		}
		exercises = append(exercises, dtoEx)
	}

	workout := &dto.Workout{
		ID:   w.Workout.ID,
		Date: w.Workout.Date,
		Name: w.Workout.Name,
		Note: w.Workout.Notes,
	}

	return &dto.WorkoutExercises{
		Workout:   workout,
		Exercises: exercises,
	}

}

func FromPaginationToFilter(pag *dto.Pagination) (*model.WorkoutsFilter, error) {
	var (
		filter model.WorkoutsFilter
		err    error
	)
	if len(pag.StartDate) != 0 {
		filter.StartDate, err = time.Parse(time.RFC3339, pag.StartDate)
		if err != nil {
			return nil, errors.New("invalid date format")
		}
	}

	if len(pag.EndDate) != 0 {
		filter.EndDate, err = time.Parse(time.RFC3339, pag.EndDate)
		if err != nil {
			return nil, errors.New("invalid date format")
		}
	}

	if len(pag.Limit) != 0 {
		filter.Limit, err = strconv.ParseUint(pag.Limit, 10, 64)
		if err != nil {
			return nil, err
		}

		if filter.Limit > 30 || filter.Limit < 1 {
			return nil, errors.New("limit must be between 1 and 30")
		}
	} else {
		filter.Limit = 10
	}

	var page uint64
	if len(pag.Page) != 0 {
		page, err = strconv.ParseUint(pag.Page, 10, 64)
		if err != nil {
			return nil, err
		}
		if page < 1 {
			return nil, errors.New("page must be greater or equal than 1")
		}
	} else {
		page = 1
	}

	filter.Offset = (page - 1) * filter.Limit
	return &filter, nil
}

func ToWorkoutResp(w *model.Workout) *dto.Workout {
	return &dto.Workout{
		ID:     w.ID,
		Date:   w.Date,
		Name:   w.Name,
		UserId: w.UserID,
		Note:   w.Notes,
	}
}

func ToWorkoutsResp(workouts []*model.Workout) []*dto.Workout {
	var wrks []*dto.Workout
	for _, w := range workouts {
		wrks = append(wrks, ToWorkoutResp(w))
	}

	return wrks
}

func FromAddExerciseToWorkout(d *dto.WorkoutExercise) *model.WorkoutExercise {
	return &model.WorkoutExercise{
		WorkoutID:  d.WorkoutID,
		ExerciseID: d.ExerciseID,
		Sets:       d.Sets,
		Reps:       d.Reps,
		Weight:     d.Weight,
		Duration:   d.Duration,
		Distance:   d.Distance,
		Exercise: model.Exercise{
			ID:          d.Exercise.ID,
			Name:        d.Exercise.Name,
			Type:        d.Exercise.Type,
			MuscleGroup: d.Exercise.MuscleGroup,
			Description: d.Exercise.Description,
		},
	}
}
