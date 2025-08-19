package workout

import (
	"context"
	"fmt"

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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate workouts: %w", err)
	}

	return workouts, nil
}

func (r *repo) AddWorkoutExercise(ctx context.Context, we *model.WorkoutExercise) (int64, error) {
	query, args, err := r.qb.Insert("workout_exercises").
		Columns("workout_id", "exercise_id", "sets", "reps", "weight", "duration", "distance").
		Values(we.WorkoutID, we.ExerciseID, we.Sets, we.Reps, we.Weight, we.Duration, we.Distance).
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

func (r *repo) GetExercisesByWorkoutID(ctx context.Context, workoutID int64) ([]*model.WorkoutExercise, error) {
	query, args, err := r.qb.
		Select("we.id", "we.workout_id", "we.exercise_id", "we.sets", "we.reps", "we.weight", "we.duration", "we.distance", "e.name", "e.type", "e.muscle_group", "e.description").
		From("workout_exercises we").
		Join("exercises e ON we.exercise_id = e.id").
		Where(squirrel.Eq{"we.workout_id": workoutID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list workouts: %w", err)
	}
	defer rows.Close()

	var exercises []*model.WorkoutExercise
	for rows.Next() {
		var exercise model.WorkoutExercise
		err = rows.Scan(
			&exercise.ID,
			&exercise.WorkoutID,
			&exercise.ExerciseID,
			&exercise.Sets,
			&exercise.Reps,
			&exercise.Weight,
			&exercise.Duration,
			&exercise.Distance,
			&exercise.Exercise.Name,
			&exercise.Exercise.Type,
			&exercise.Exercise.MuscleGroup,
			&exercise.Exercise.Description,
		)

		exercises = append(exercises, &exercise)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate workouts: %w", err)
	}

	return exercises, nil
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

func (r *repo) AddRecord(ctx context.Context, user *model.UserRecord) (int64, error) {
	query, args, err := r.qb.
		Insert("personal_records").
		Columns("user_id", "exercise_id", "weight", "reps", "date").
		Values(user.UserID, user.ExerciseID, user.Weight, user.Reps, user.Date).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build insert query: %w", err)
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert record: %w", err)
	}

	return id, nil
}

func (r *repo) GetPersonalRecord(ctx context.Context, userID, exerciseID int64) (*model.UserRecord, error) {
	query, args, err := r.qb.
		Select("weight", "reps").
		From("personal_records").
		Where(squirrel.Eq{"user_id": userID, "exercise_id": exerciseID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var record model.UserRecord

	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(&record.Weight, &record.Reps)
	if err != nil {
		return nil, fmt.Errorf("failed to list workouts: %w", err)
	}

	return &record, nil
}

func (r *repo) UpdatePersonalRecord(ctx context.Context, user *model.UserRecord) error {
	query, args, err := r.qb.
		Update("personal_records").
		Set("weight", user.Weight).
		Set("reps", user.Reps).
		Set("date", user.Date).
		Where(squirrel.Eq{"user_id": user.UserID, "exercise_id": user.ExerciseID}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = r.db.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update records: %w", err)
	}
	return nil
}

func (r *repo) ListRecords(ctx context.Context, userId int64) ([]*model.UserRecord, error) {
	query, args, err := r.qb.
		Select("pr.id", "pr.exercise_id", "pr.weight", "pr.reps", "pr.date", "e.name", "e.type", "e.muscle_group", "e.description").
		From("personal_records pr").
		Join("exercises e ON pr.exercise_id = e.id").
		Where(squirrel.Eq{"user_id": userId}).OrderBy("pr.date DESC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list records: %w", err)
	}
	defer rows.Close()

	var records []*model.UserRecord
	for rows.Next() {
		var record model.UserRecord
		err = rows.Scan(&record.ID, &record.ExerciseID, &record.Weight, &record.Reps, &record.Date, &record.Exercise.Name, &record.Exercise.Type, &record.Exercise.MuscleGroup, &record.Exercise.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan records: %w", err)
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list records: %w", err)
	}
	return records, nil
}
