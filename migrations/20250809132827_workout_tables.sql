-- +goose Up
-- +goose StatementBegin
create table users (
    id int generated always as identity primary key,
    name VARCHAR(50) UNIQUE NOT NULL,
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
    record_type VARCHAR(20) NOT NULL CHECK (record_type IN ('weight', 'reps', 'sets', 'distance', 'duration')),
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
    sets INTEGER DEFAULT 1,
    reps INTEGER DEFAULT 0,
    weight DECIMAL(5,2) DEFAULT 0,
    duration INTEGER DEFAULT 0, -- для кардио в секундах
    distance DECIMAL(6,2) DEFAULT 0, -- для кардио в км
    notes TEXT,
    created_at timestamp not null default now()
);

CREATE TABLE IF NOT EXISTS personal_records (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    exercise_id INTEGER REFERENCES exercises(id),
    weight DECIMAL(5,2) DEFAULT NULL,
    reps INTEGER DEFAULT NULL,
    sets INTEGER DEFAULT NULL,
    duration INTEGER DEFAULT NULL,
    distance DECIMAL(7,2) DEFAULT NULL,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    UNIQUE(user_id, exercise_id)
);

CREATE INDEX IF NOT EXISTS idx_workouts_user_id ON workouts(user_id);
CREATE INDEX IF NOT EXISTS idx_workouts_date ON workouts(date);
CREATE INDEX IF NOT EXISTS idx_workout_exercises_workout_id ON workout_exercises(workout_id);
CREATE INDEX IF NOT EXISTS idx_workout_exercises_exercise_id ON workout_exercises(exercise_id);
CREATE INDEX IF NOT EXISTS idx_personal_records_user_id ON personal_records(user_id);
CREATE INDEX IF NOT EXISTS idx_personal_records_exercise_id ON personal_records(exercise_id);

INSERT INTO exercises (name, type, muscle_group, description, record_type) VALUES
-- Силовые упражнения
('Жим лежа', 'strength', 'Грудь', 'Базовое упражнение для развития грудных мышц, передних дельт и трицепсов', 'weight'),
('Приседания со штангой', 'strength', 'Ноги', 'Базовое упражнение для развития квадрицепсов, ягодичных мышц и задней поверхности бедра', 'weight'),
('Становая тяга', 'strength', 'Спина', 'Базовое упражнение для развития мышц спины, ног и укрепления всего тела', 'weight'),
('Подтягивания', 'strength', 'Спина', 'Упражнение для развития широчайших мышц спины и бицепсов', 'reps'),
('Отжимания', 'strength', 'Грудь', 'Упражнение с собственным весом для развития грудных мышц, трицепсов и дельт', 'reps'),
('Жим штанги стоя', 'strength', 'Плечи', 'Упражнение для развития дельтовидных мышц и стабилизаторов корпуса', 'weight'),
('Тяга штанги в наклоне', 'strength', 'Спина', 'Упражнение для развития широчайших мышц спины и задних дельт', 'weight'),
('Сгибание рук со штангой', 'strength', 'Бицепс', 'Изолирующее упражнение для развития бицепсов', 'weight'),
('Французский жим', 'strength', 'Трицепс', 'Изолирующее упражнение для развития трицепсов', 'weight'),
('Подъемы на носки', 'strength', 'Голени', 'Упражнение для развития икроножных мышц', 'weight'),
('Планка', 'strength', 'Пресс', 'Статическое упражнение для укрепления мышц кора', 'duration'),
('Скручивания', 'strength', 'Пресс', 'Упражнение для развития прямых мышц живота', 'reps'),

-- Кардио упражнения
('Бег', 'cardio', 'Кардио', 'Кардиотренировка для развития выносливости и сжигания калорий', 'distance'),
('Быстрая ходьба', 'cardio', 'Кардио', 'Низкоинтенсивная кардиотренировка подходящая для начинающих', 'duration'),
('Велосипед', 'cardio', 'Кардио', 'Кардиотренировка на велосипеде или велотренажере', 'distance'),
('Эллиптический тренажер', 'cardio', 'Кардио', 'Кардиотренировка на эллиптическом тренажере', 'distance'),
('Плавание', 'cardio', 'Кардио', 'Комплексная кардиотренировка задействующая все группы мышц', 'distance'),
('Гребля', 'cardio', 'Кардио', 'Кардиотренировка на гребном тренажере', 'distance'),
('Степпер', 'cardio', 'Кардио', 'Кардиотренировка имитирующая подъем по лестнице', 'duration'),
('Прыжки на скакалке', 'cardio', 'Кардио', 'Высокоинтенсивная кардиотренировка для развития координации', 'sets'),
('HIIT тренировка', 'cardio', 'Кардио', 'Высокоинтенсивная интервальная тренировка', 'duration'),
('Танцы', 'cardio', 'Кардио', 'Кардиотренировка в виде танцевальных движений', 'duration')
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
