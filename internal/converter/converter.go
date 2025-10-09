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
	}
}

func FromUserLoginRequest(u *dto.UserLoginRequest) *model.LoginUserParams {
	return &model.LoginUserParams{
		Email:    u.Email,
		Password: u.Password,
	}
}

func FromCreateWorkoutRequest(r *dto.CreateWorkoutRequest) (*model.Workout, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, r.WorkoutDate)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse workout date")
	}
	return model.NewWorkout(r.UserID, r.Name, r.Description, t), nil
}

func ToGetWorkoutResp(w *model.Workout) *dto.Workout {
	if w == nil {
		return nil
	}

	// Конвертация упражнения
	convertExercise := func(ex *model.WorkoutExercise) *dto.WorkoutExercise {
		if ex == nil {
			return nil
		}

		// Преобразование подходов
		var sets []*dto.ExerciseSet
		for _, s := range ex.Sets {
			sets = append(sets, &dto.ExerciseSet{
				ID:        s.ID,
				SetNumber: s.SetNumber,
				Weight:    s.Weight,
				Reps:      s.Reps,
			})
		}

		var cardio *dto.CardioRecord
		if ex.Cardio != nil {
			cardio = &dto.CardioRecord{
				ID:              ex.Cardio.ID,
				DistanceKm:      ex.Cardio.DistanceKm,
				DurationSeconds: ex.Cardio.DurationSeconds,
			}
		}

		// Преобразование связанного упражнения
		var exerciseDto *dto.Exercise
		if ex.Exercise != nil {
			exerciseDto = &dto.Exercise{
				ID:          ex.Exercise.ID,
				Name:        ex.Exercise.Name,
				Description: ex.Exercise.Description,
				Type:        string(ex.Exercise.Type),
				MuscleGroup: ex.Exercise.MuscleGroup,
			}
		}

		return &dto.WorkoutExercise{
			ID:         ex.ID,
			WorkoutID:  ex.WorkoutID,
			ExerciseID: ex.ExerciseID,
			Notes:      ex.Notes,
			Exercise:   exerciseDto,
			Sets:       sets,
			Cardio:     cardio,
		}
	}

	// Преобразование всех упражнений тренировки
	var exercises []*dto.WorkoutExercise
	for _, exercise := range w.Exercises {
		exercises = append(exercises, convertExercise(exercise))
	}

	return &dto.Workout{
		ID:        w.ID,
		UserID:    w.UserID,
		Name:      w.Name,
		Notes:     w.Notes,
		Date:      w.Date,
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

func ToWorkoutsResp(workouts []*model.Workout) []*dto.Workout {
	if workouts == nil {
		return nil
	}

	result := make([]*dto.Workout, 0, len(workouts))
	for _, w := range workouts {
		result = append(result, ToGetWorkoutResp(w))
	}
	return result
}

func FromAddExerciseToWorkout(d *dto.AddExerciseRequest, workoutID int64) *model.WorkoutExercise {
	we := model.NewWorkoutExercise(workoutID, d.ExerciseID, d.Notes)
	for _, set := range d.Sets {
		s := model.NewExerciseSet(0, set.SetNumber, set.Weight, set.Reps)
		we.Sets = append(we.Sets, s)
	}
	return we
}

func ToListExercisesResp(exercises []*model.Exercise) []*dto.Exercise {
	var wrks []*dto.Exercise
	for _, ex := range exercises {
		e := &dto.Exercise{
			ID:          ex.ID,
			Name:        ex.Name,
			Type:        ex.Type,
			MuscleGroup: ex.MuscleGroup,
			Description: ex.Description,
		}

		wrks = append(wrks, e)
	}

	return wrks
}

func ToPersonalRecord(records []*model.PersonalRecord) []*dto.PersonalRecord {
	if records == nil {
		return nil
	}

	result := make([]*dto.PersonalRecord, 0, len(records))
	for _, r := range records {
		result = append(result, toPersonalRecordDTO(r))
	}
	return result
}

func toPersonalRecordDTO(r *model.PersonalRecord) *dto.PersonalRecord {
	if r == nil {
		return nil
	}

	return &dto.PersonalRecord{
		ID:                r.ID,
		UserID:            r.UserID,
		ExerciseID:        r.ExerciseID,
		RecordType:        string(r.RecordType),
		Value:             r.Value,
		WorkoutExerciseID: r.WorkoutExerciseID,
		AchievedAt:        r.AchievedAt,
		Exercise:          toExerciseDTO(r.Exercise),
	}
}

func toExerciseDTO(exercise *model.Exercise) *dto.Exercise {
	return &dto.Exercise{
		ID:          exercise.ID,
		Name:        exercise.Name,
		Description: exercise.Description,
		Type:        string(exercise.Type),
		MuscleGroup: exercise.MuscleGroup,
	}
}

func FromUpdateWorkoutRequest(d *dto.UpdateWorkoutRequest) *model.UpdateWorkout {
	var wD *time.Time
	if d.WorkoutDate != nil && *d.WorkoutDate != "" {
		layout := "2006-01-02" // Формат YYYY-MM-DD
		parsedDate, _ := time.Parse(layout, *d.WorkoutDate)
		wD = &parsedDate
	}
	return &model.UpdateWorkout{
		Name:  d.Name,
		Notes: d.Description,
		Date:  wD,
	}
}

func FromAddCardioToWorkoutRequest(d *dto.AddCardioRequest) *model.CardioRecord {
	return model.NewCardioRecord(0, d.DurationSeconds, d.DistanceKm)
}

func ToCreateWorkoutsResp(resp *model.Workout) *dto.Workout {
	return &dto.Workout{
		ID:     resp.ID,
		UserID: resp.UserID,
		Name:   resp.Name,
		Notes:  resp.Notes,
		Date:   resp.Date,
	}
}
