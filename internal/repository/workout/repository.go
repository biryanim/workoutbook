package workout

import (
	"context"
	"fmt"
	"time"

	apperrors "github.com/biryanim/workoutbook/internal/errors"

	"github.com/Masterminds/squirrel"
	"github.com/biryanim/workoutbook/internal/client/db"
	"github.com/biryanim/workoutbook/internal/model"
	"github.com/biryanim/workoutbook/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

var _ repository.WorkoutRepository = (*repo)(nil)

type repo struct {
	db db.Client
	qb squirrel.StatementBuilderType
}

func NewRepository(db db.Client) *repo {
	return &repo{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *repo) CreateWorkout(ctx context.Context, workout *model.Workout) (int64, error) {
	query, args, err := r.qb.
		Insert("workouts").
		Columns("user_id", "date", "name", "notes").
		Values(workout.UserID, workout.Date, workout.Name, workout.Notes).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build insert query: %w", err)
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert workout: %w", err)
	}

	return id, nil
}

func (r *repo) GetWorkoutByID(ctx context.Context, workoutID, userId int64) (*model.Workout, error) {
	query, args, err := r.qb.
		Select("id", "user_id", "date", "notes", "name", "created_at", "updated_at").
		From("workouts").
		Where(squirrel.Eq{"id": workoutID, "user_id": userId}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var workout model.Workout
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(
		&workout.ID,
		&workout.UserID,
		&workout.Date,
		&workout.Notes,
		&workout.Name,
		&workout.CreatedAt,
		&workout.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrTaskNotFound
		}
	}
	return &workout, nil
}

func (r *repo) ListWorkoutsByUserID(ctx context.Context, userId int64) ([]*model.Workout, error) {
	query, args, err := r.qb.
		Select("id", "user_id", "name", "notes", "date", "created_at", "updated_at").
		From("workouts").
		Where(squirrel.Eq{"user_id": userId}).
		OrderBy("date DESC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list workouts: %w", err)
	}
	defer rows.Close()

	var workouts []*model.Workout
	for rows.Next() {
		var v model.Workout
		err = rows.Scan(
			&v.ID,
			&v.UserID,
			&v.Name,
			&v.Notes,
			&v.Date,
			&v.CreatedAt,
			&v.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workouts: %w", err)
		}
		workouts = append(workouts, &v)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return workouts, nil
}

func (r *repo) ListWorkouts(ctx context.Context, userId int64, filter *model.WorkoutsFilter) ([]*model.Workout, error) {
	builder := r.qb.Select("id", "user_id", "date", "notes", "name", "created_at", "updated_at").
		From("workouts").
		Where(squirrel.Eq{"user_id": userId}).
		OrderBy(
			"date DESC",
		).
		Limit(filter.Limit).
		Offset(filter.Offset)

	if !filter.StartDate.IsZero() {
		builder = builder.Where(squirrel.GtOrEq{"created_at": filter.StartDate})
	}

	if !filter.EndDate.IsZero() {
		builder = builder.Where(squirrel.LtOrEq{"created_at": filter.EndDate})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var workouts []*model.Workout
	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list workouts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var workout model.Workout
		err = rows.Scan(
			&workout.ID,
			&workout.UserID,
			&workout.Date,
			&workout.Notes,
			&workout.Name,
			&workout.CreatedAt,
			&workout.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workout: %w", err)
		}
		workouts = append(workouts, &workout)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate workouts: %w", err)
	}

	return workouts, nil
}

func (r *repo) UpdateWorkout(ctx context.Context, workout *model.Workout) error {
	fmt.Println(workout.Date)
	query, args, err := r.qb.
		Update("workouts").
		Set("name", workout.Name).
		Set("notes", workout.Notes).
		Set("date", workout.Date).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": workout.ID}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}
	_, err = r.db.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update workout: %w", err)
	}

	return nil
}

func (r *repo) DeleteWorkout(ctx context.Context, id int64) error {
	query, args, err := r.qb.
		Delete("workouts").
		Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = r.db.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete workout: %w", err)
	}
	return nil
}

func (r *repo) AddWorkoutExercise(ctx context.Context, we *model.WorkoutExercise) error {
	query, args, err := r.qb.
		Insert("workout_exercises").
		Columns("workout_id", "exercise_id", "notes").
		Values(we.WorkoutID, we.ExerciseID, we.Notes).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(&we.ID)
	if err != nil {
		return fmt.Errorf("failed to add workout exercise: %w", err)
	}

	return nil
}

func (r *repo) GetWorkoutExercises(ctx context.Context, workoutID int64) ([]*model.WorkoutExercise, error) {
	query, args, err := r.qb.
		Select("we.id", "we.workout_id", "we.exercise_id", "we.notes", "we.created_at",
			"e.id", "e.name", "e.description", "e.type", "e.muscle_group", "e.created_at").
		From("workout_exercises we").
		Join("exercises e ON we.exercise_id = e.id").
		Where(squirrel.Eq{"we.workout_id": workoutID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout exercises: %w", err)
	}
	defer rows.Close()

	var exercises []*model.WorkoutExercise
	for rows.Next() {
		var we model.WorkoutExercise
		var ex model.Exercise
		err = rows.Scan(
			&we.ID,
			&we.WorkoutID,
			&we.ExerciseID,
			&we.Notes,
			&we.CreatedAt,
			&ex.ID,
			&ex.Name,
			&ex.Description,
			&ex.Type,
			&ex.MuscleGroup,
			&ex.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workout exercises: %w", err)
		}
		we.Exercise = &ex
		exercises = append(exercises, &we)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan workout exercises: %w", err)
	}
	return exercises, nil
}

func (r *repo) DeleteWorkoutExercise(ctx context.Context, workoutExerciseID int64) error {
	query, args, err := r.qb.
		Delete("workout_exercises").
		Where(squirrel.Eq{"id": workoutExerciseID}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}
	_, err = r.db.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete workout exercise: %w", err)
	}

	return nil
}

func (r *repo) AddExerciseSet(ctx context.Context, set *model.ExerciseSet) error {
	query, args, err := r.qb.
		Insert("exercise_sets").
		Columns("workout_exercise_id", "set_number", "weight", "reps").
		Values(set.WorkoutExerciseID, set.SetNumber, set.Weight, set.Reps).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(&set.ID)
	if err != nil {
		return fmt.Errorf("failed to add workout exercise_sets: %w", err)
	}

	return nil
}

func (r *repo) GetExerciseSets(ctx context.Context, workoutExerciseID int64) ([]*model.ExerciseSet, error) {
	query, args, err := r.qb.
		Select("id", "workout_exercise_id", "set_number", "weight", "reps", "created_at").
		From("exercise_sets").
		Where(squirrel.Eq{"workout_exercise_id": workoutExerciseID}).OrderBy("set_number").ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout exercise_sets: %w", err)
	}
	defer rows.Close()

	var sets []*model.ExerciseSet
	for rows.Next() {
		var set model.ExerciseSet
		err = rows.Scan(
			&set.ID,
			&set.WorkoutExerciseID,
			&set.SetNumber,
			&set.Weight,
			&set.Reps,
			&set.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workout exercise_sets: %w", err)
		}
		sets = append(sets, &set)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get workout exercise_sets: %w", err)
	}
	return sets, nil
}

func (r *repo) IsUserHaveWorkout(ctx context.Context, userId, workoutId int64) (bool, error) {
	query, args, err := r.qb.
		Select("count(*)").
		From("workouts").
		Where(squirrel.Eq{"id": workoutId, "user_id": userId}).ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build select query: %w", err)
	}

	var count int
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if workout is have: %w", err)
	}

	return count > 0, nil
}

func (r *repo) GetExercises(ctx context.Context, exerciseType string) ([]*model.Exercise, error) {
	builder := r.qb.Select("id", "name", "type", "muscle_group", "description").
		From("exercises")
	if exerciseType != "" {
		builder = builder.Where("type = ?", exerciseType)
	}
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list workouts: %w", err)
	}
	defer rows.Close()

	var exercises []*model.Exercise

	for rows.Next() {
		var exercise model.Exercise
		err = rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Type,
			&exercise.MuscleGroup,
			&exercise.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workout: %w", err)
		}

		exercises = append(exercises, &exercise)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate workouts: %w", err)
	}

	return exercises, nil
}

func (r *repo) GetExerciseSetByID(ctx context.Context, id int64) (*model.ExerciseSet, error) {
	query, args, err := r.qb.
		Select("id", "workout_exercise_id", "set_number", "weight", "reps", "created_at").
		From("exercise_sets").
		Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var set model.ExerciseSet
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(
		&set.ID,
		&set.WorkoutExerciseID,
		&set.SetNumber,
		&set.Weight,
		&set.Reps,
		&set.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout exercise_sets: %w", err)
	}
	return &set, nil
}

func (r *repo) UpdateExerciseSet(ctx context.Context, set *model.ExerciseSet) error {
	query, args, err := r.qb.
		Update("exercise_sets").
		Set("weight", set.Weight).
		Set("reps", set.Reps).
		Where(squirrel.Eq{"id": set.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update exercise_sets: %w", err)
	}

	_, err = r.db.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update exercise_sets: %w", err)
	}
	return nil
}

func (r *repo) DeleteExerciseSet(ctx context.Context, id int64) error {
	query, args, err := r.qb.
		Delete("exercise_sets").
		Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete exercise_sets: %w", err)
	}

	_, err = r.db.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete exercise_sets: %w", err)
	}
	return nil
}

func (r *repo) ReorderExerciseSets(ctx context.Context, workoutExerciseID int64) error {
	query := `
        WITH numbered_sets AS (
            SELECT id, ROW_NUMBER() OVER (ORDER BY set_number) as new_number
            FROM exercise_sets
            WHERE workout_exercise_id = $1
        )
        UPDATE exercise_sets
        SET set_number = numbered_sets.new_number
        FROM numbered_sets
        WHERE exercise_sets.id = numbered_sets.id
    `

	_, err := r.db.DB().ExecContext(ctx, query, workoutExerciseID)
	if err != nil {
		return fmt.Errorf("failed to reorder sets: %w", err)
	}

	return nil
}

func (r *repo) GetPersonalRecord(ctx context.Context, userID, exerciseID int64, recordType model.RecordType) (*model.PersonalRecord, error) {
	query, args, err := r.qb.
		Select("id", "user_id", "exercise_id", "record_type", "value", "workout_exercise_id", "achieved_at").
		From("personal_records").
		Where(squirrel.Eq{"user_id": userID, "exercise_id": exerciseID, "record_type": recordType}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var personalRecord model.PersonalRecord
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(
		&personalRecord.ID,
		&personalRecord.UserID,
		&personalRecord.ExerciseID,
		&personalRecord.RecordType,
		&personalRecord.Value,
		&personalRecord.WorkoutExerciseID,
		&personalRecord.AchievedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan personal record: %w", err)
	}
	return &personalRecord, nil
}

func (r *repo) GetPersonalRecords(ctx context.Context, userID int64) ([]*model.PersonalRecord, error) {
	query, args, err := r.qb.
		Select("pr.id", "pr.user_id", "pr.exercise_id", "pr.record_type", "pr.value", "pr.workout_exercise_id", "pr.achieved_at",
			"e.id", "e.name", "e.description", "e.type", "e.muscle_group", "e.created_at").
		From("personal_records pr").
		Join("exercises e ON pr.exercise_id = e.id").
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("pr.achieved_at DESC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get personal records: %w", err)
	}
	defer rows.Close()

	var records []*model.PersonalRecord
	for rows.Next() {
		var pr model.PersonalRecord
		var ex model.Exercise

		err = rows.Scan(
			&pr.ID,
			&pr.UserID,
			&pr.ExerciseID,
			&pr.RecordType,
			&pr.Value,
			&pr.WorkoutExerciseID,
			&pr.AchievedAt,
			&ex.ID,
			&ex.Name,
			&ex.Description,
			&ex.Type,
			&ex.MuscleGroup,
			&ex.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan personal records: %w", err)
		}
		pr.Exercise = &ex
		records = append(records, &pr)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get personal records: %w", err)
	}
	return records, nil
}

func (r *repo) UpsertPersonalRecord(ctx context.Context, record *model.PersonalRecord) error {
	query, args, err := r.qb.
		Insert("personal_records").
		Columns("user_id", "exercise_id", "record_type", "value", "workout_exercise_id", "achieved_at").
		Values(record.UserID, record.ExerciseID, record.RecordType, record.Value, record.WorkoutExerciseID, record.AchievedAt).
		Suffix(`
			on conflict (user_id, exercise_id, record_type)
			do update set
				value = EXCLUDED.value,
				workout_exercise_id = EXCLUDED.workout_exercise_id,
				achieved_at = EXCLUDED.achieved_at
			RETURNING id
		`).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update personal records: %w", err)
	}

	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(&record.ID)
	if err != nil {
		return fmt.Errorf("failed to upsert personal records: %w", err)
	}

	return nil
}

func (r *repo) AddCardioRecord(ctx context.Context, cardio *model.CardioRecord) error {
	query, args, err := r.qb.
		Insert("cardio_records").
		Columns("workout_exercise_id", "distance_km", "duration_seconds").
		Values(cardio.WorkoutExerciseID, cardio.DistanceKm, cardio.DurationSeconds).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update cardio records: %w", err)
	}

	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(&cardio.ID)
	if err != nil {
		return fmt.Errorf("failed to update cardio records: %w", err)
	}

	return nil
}

func (r *repo) GetCardioRecord(ctx context.Context, workoutExerciseID int64) (*model.CardioRecord, error) {
	query, args, err := r.qb.
		Select("id", "workout_exercise_id", "distance_km", "duration_seconds", "created_at").
		From("cardio_records").
		Where(squirrel.Eq{"workout_exercise_id": workoutExerciseID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var card model.CardioRecord
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(
		&card.ID,
		&card.WorkoutExerciseID,
		&card.DistanceKm,
		&card.DurationSeconds,
		&card.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan cardio record: %w", err)
	}

	return &card, nil
}

func (r *repo) UpdateCardioRecord(ctx context.Context, cardio *model.CardioRecord) error {
	query, args, err := r.qb.
		Update("cardio_records").
		Set("distance_km", cardio.DistanceKm).
		Set("duration_seconds", cardio.DurationSeconds).
		Where(squirrel.Eq{"id": cardio.ID}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update cardio records: %w", err)
	}

	_, err = r.db.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update cardio records: %w", err)
	}

	return nil
}

func (r *repo) DeleteCardioRecord(ctx context.Context, id int64) error {
	query, args, err := r.qb.
		Delete("cardio_records").
		Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete cardio records: %w", err)
	}
	_, err = r.db.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete cardio records: %w", err)
	}
	return nil
}

func (r *repo) GetExerciseByID(ctx context.Context, exerciseID int64) (*model.Exercise, error) {
	query, args, err := r.qb.
		Select("id", "name", "type", "muscle_group", "description", "created_at").
		From("exercises").
		Where(squirrel.Eq{"id": exerciseID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var exercise model.Exercise
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Type,
		&exercise.MuscleGroup,
		&exercise.Description,
		&exercise.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get exercise by id: %w", err)
	}

	return &exercise, nil
}

func (r *repo) CreateExerciseSet(ctx context.Context, workoutExerciseID int64, setNumber int, weight float64, reps int) error {
	query, args, err := r.qb.
		Insert("exercise_sets").
		Columns("workout_exercise_id", "set_number", "weight", "reps").
		Values(workoutExerciseID, setNumber, weight, reps).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert exercise_sets: %w", err)
	}

	_, err = r.db.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert exercise_sets: %w", err)
	}
	return nil
}
