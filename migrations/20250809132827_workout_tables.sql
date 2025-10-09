-- +goose Up
-- +goose StatementBegin
create table users (
    id int generated always as identity primary key,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at timestamp not null default now(),
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS exercises (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('strength', 'cardio')),
    muscle_group VARCHAR(50),
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS workouts (
    id int generated always as identity primary key,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(100) NOT NULL,
    notes TEXT,
    created_at timestamp not null default now(),
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS workout_exercises (
    id int generated always as identity primary key,
    workout_id INTEGER REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id INTEGER REFERENCES exercises(id),
    notes TEXT,
    created_at timestamp not null default now()
);

create table exercise_sets(
    id int generated always as identity primary key,
    workout_exercise_id int not null references workout_exercises(id) on delete cascade ,
    set_number int not null,
    weight decimal(6, 2),
    reps int,
    created_at timestamp default now(),
    unique (workout_exercise_id, set_number)
);

create table cardio_records(
    id int generated always as identity primary key ,
    workout_exercise_id int not null references workout_exercises(id) on delete cascade ,
    distance_km decimal(6,2),
    duration_seconds int,
    created_at timestamp default now()
);

CREATE TABLE IF NOT EXISTS personal_records (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    exercise_id INTEGER REFERENCES exercises(id),
    record_type VARCHAR(20) NOT NULL CHECK (record_type IN ('max_weight', 'max_reps', 'max_distance', 'best_time')),
    value DECIMAL(10, 2) NOT NULL,
    workout_exercise_id INTEGER REFERENCES workout_exercises(id) ON DELETE SET NULL,
    achieved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, exercise_id, record_type)
);

CREATE INDEX IF NOT EXISTS idx_workouts_user_id ON workouts(user_id);
CREATE INDEX IF NOT EXISTS idx_workouts_date ON workouts(date);
CREATE INDEX IF NOT EXISTS idx_workout_exercises_workout_id ON workout_exercises(workout_id);
CREATE INDEX IF NOT EXISTS idx_workout_exercises_exercise_id ON workout_exercises(exercise_id);
CREATE INDEX IF NOT EXISTS idx_personal_records_user_id ON personal_records(user_id);
CREATE INDEX IF NOT EXISTS idx_personal_records_exercise_id ON personal_records(exercise_id);

INSERT INTO exercises (name, type, muscle_group, description) VALUES
    ('Жим лежа', 'strength', 'Грудь', 'Базовое упражнение для развития грудных мышц, передних дельт и трицепсов'),
    ('Приседания со штангой', 'strength', 'Ноги', 'Базовое упражнение для развития квадрицепсов, ягодичных мышц и задней поверхности бедра'),
    ('Становая тяга', 'strength', 'Спина', 'Базовое упражнение для развития мышц спины, ног и укрепления всего тела'),
    ('Подтягивания', 'strength', 'Спина', 'Упражнение для развития широчайших мышц спины и бицепсов'),
    ('Отжимания', 'strength', 'Грудь', 'Упражнение с собственным весом для развития грудных мышц, трицепсов и дельт'),
    ('Жим штанги стоя', 'strength', 'Плечи', 'Упражнение для развития дельтовидных мышц и стабилизаторов корпуса'),
    ('Тяга штанги в наклоне', 'strength', 'Спина', 'Упражнение для развития широчайших мышц спины и задних дельт'),
    ('Сгибание рук со штангой', 'strength', 'Бицепс', 'Изолирующее упражнение для развития бицепсов'),
    ('Французский жим', 'strength', 'Трицепс', 'Изолирующее упражнение для развития трицепсов'),
    ('Подъемы на носки', 'strength', 'Голени', 'Упражнение для развития икроножных мышц'),
    ('Планка', 'strength', 'Пресс', 'Статическое упражнение для укрепления мышц кора'),
    ('Скручивания', 'strength', 'Пресс', 'Упражнение для развития прямых мышц живота'),
    ('Бег', 'cardio', 'Кардио', 'Кардиотренировка для развития выносливости и сжигания калорий'),
    ('Быстрая ходьба', 'cardio', 'Кардио', 'Низкоинтенсивная кардиотренировка подходящая для начинающих'),
    ('Велосипед', 'cardio', 'Кардио', 'Кардиотренировка на велосипеде или велотренажере'),
    ('Эллиптический тренажер', 'cardio', 'Кардио', 'Кардиотренировка на эллиптическом тренажере'),
    ('Плавание', 'cardio', 'Кардио', 'Комплексная кардиотренировка задействующая все группы мышц'),
    ('Гребля', 'cardio', 'Кардио', 'Кардиотренировка на гребном тренажере'),
    ('Степпер', 'cardio', 'Кардио', 'Кардиотренировка имитирующая подъем по лестнице'),
    ('Прыжки на скакалке', 'cardio', 'Кардио', 'Высокоинтенсивная кардиотренировка для развития координации'),
    ('HIIT тренировка', 'cardio', 'Кардио', 'Высокоинтенсивная интервальная тренировка'),
    ('Танцы', 'cardio', 'Кардио', 'Кардиотренировка в виде танцевальных движений')
ON CONFLICT (name) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_personal_records_exercise_id;
DROP INDEX IF EXISTS idx_personal_records_user_id;
DROP INDEX IF EXISTS idx_workout_exercises_exercise_id;
DROP INDEX IF EXISTS idx_workout_exercises_workout_id;
DROP INDEX IF EXISTS idx_workouts_date;
DROP INDEX IF EXISTS idx_workouts_user_id;

DROP TABLE IF EXISTS personal_records CASCADE;
DROP TABLE IF EXISTS workout_exercises CASCADE;
DROP TABLE IF EXISTS workouts CASCADE;
DROP TABLE IF EXISTS exercises CASCADE;
DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd
