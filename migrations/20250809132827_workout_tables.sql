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
    id int generated always as identity primary key,
    name VARCHAR(100) NOT NULL unique ,
    type VARCHAR(20) NOT NULL CHECK (type IN ('strength', 'cardio')),
    muscle_group VARCHAR(50),
    description TEXT,
    created_at timestamp not null default now()
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
    id int generated always as identity primary key,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    exercise_id INTEGER REFERENCES exercises(id),
    weight DECIMAL(5,2),
    reps INTEGER,
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

INSERT INTO exercises (name, type, muscle_group, description) VALUES
-- Силовые упражнения
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

-- Кардио упражнения
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
