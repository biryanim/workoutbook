// Загрузка и отображение списка тренировок
async function loadWorkouts() {
    const container = document.getElementById('workouts-list');

    try {
        const data = await WorkoutAPI.getWorkouts();
        const workouts = data.workouts || [];

        if (workouts.length === 0) {
            container.innerHTML = `
                <div style="text-align: center; padding: 60px 20px; color: var(--text-muted);">
                    <p style="font-size: 18px; margin-bottom: 20px;">У вас пока нет тренировок</p>
                    <p>Создайте свою первую тренировку, чтобы начать отслеживать прогресс</p>
                </div>
            `;
            return;
        }

        container.innerHTML = workouts.map(workout => `
            <div class="workout-card">
                <div onclick="openWorkout(${workout.id})" style="flex: 1; cursor: pointer;">
                    <h3>${escapeHtml(workout.name)}</h3>
                    <p class="date">${formatDate(workout.workout_date)}</p>
                    ${workout.description ? `<p class="description">${escapeHtml(workout.description)}</p>` : ''}
                </div>
                <div class="workout-actions">
                    <button class="btn btn-secondary btn-sm edit-workout" data-workout-id="${workout.id}" 
                            data-workout-name="${escapeHtml(workout.name)}"
                            data-workout-description="${escapeHtml(workout.description || '')}"
                            data-workout-date="${workout.workout_date}">
                        ✏️ Изменить
                    </button>
                    <button class="btn btn-danger btn-sm delete-workout" data-workout-id="${workout.id}">
                        🗑️ Удалить
                    </button>
                </div>
            </div>
        `).join('');

        // Навешиваем обработчики на кнопки
        attachWorkoutActionHandlers();
    } catch (error) {
        container.innerHTML = `<div class="error-message">${error.message}</div>`;
    }
}

function attachWorkoutActionHandlers() {
    // Удаление тренировки
    document.querySelectorAll('.delete-workout').forEach(btn => {
        btn.addEventListener('click', async (e) => {
            e.stopPropagation();
            const workoutId = btn.getAttribute('data-workout-id');

            if (confirm('Вы уверены, что хотите удалить эту тренировку? Все упражнения будут удалены.')) {
                try {
                    await WorkoutAPI.deleteWorkout(workoutId);
                    loadWorkouts(); // Перезагружаем список
                } catch (error) {
                    alert('Ошибка: ' + error.message);
                }
            }
        });
    });

    // Редактирование тренировки
    document.querySelectorAll('.edit-workout').forEach(btn => {
        btn.addEventListener('click', (e) => {
            e.stopPropagation();
            const workoutId = btn.getAttribute('data-workout-id');
            const workoutName = btn.getAttribute('data-workout-name');
            const workoutDescription = btn.getAttribute('data-workout-description');
            const workoutDate = btn.getAttribute('data-workout-date');

            openEditWorkoutModal(workoutId, workoutName, workoutDescription, workoutDate);
        });
    });
}

function openEditWorkoutModal(workoutId, name, description, date) {
    const modal = document.getElementById('editWorkoutModal');

    document.getElementById('editWorkoutId').value = workoutId;
    document.getElementById('editWorkoutName').value = name;
    document.getElementById('editWorkoutDescription').value = description;
    document.getElementById('editWorkoutDate').value = date.split('T')[0]; // Формат YYYY-MM-DD

    modal.style.display = 'flex';
}


// Открыть детали тренировки
function openWorkout(id) {
    window.location.href = `/workout-detail.html?id=${id}`;
}

// Загрузка деталей тренировки
async function loadWorkoutDetail(workoutId) {
    const container = document.getElementById('exercises-list');

    try {
        const data = await WorkoutAPI.getWorkout(workoutId);
        const workout = data.workout;

        // Обновляем заголовок
        document.getElementById('workoutName').textContent = workout.name;
        document.getElementById('workoutDate').textContent = formatDate(workout.workout_date);

        const exercises = workout.exercises || [];

        if (exercises.length === 0) {
            container.innerHTML = `
                <div style="text-align: center; padding: 40px 20px; color: var(--text-muted);">
                    <p>В этой тренировке пока нет упражнений</p>
                </div>
            `;
            return;
        }

        container.innerHTML = exercises.map(ex => {
            if (ex.exercise.type === 'strength') {
                return renderStrengthExercise(ex);
            } else {
                return renderCardioExercise(ex);
            }
        }).join('');

        // Добавляем обработчики удаления
        attachDeleteHandlers();
    } catch (error) {
        container.innerHTML = `<div class="error-message">${error.message}</div>`;
    }
}

// Отрисовка силового упражнения
function renderStrengthExercise(ex) {
    const sets = ex.sets || [];

    return `
        <div class="exercise-card">
            <div class="exercise-header">
                <div>
                    <strong>${escapeHtml(ex.exercise.name)}</strong>
                    <span class="exercise-type strength">Силовое</span>
                </div>
            </div>
            ${ex.notes ? `<p style="color: var(--text-muted); margin-bottom: 10px;">${escapeHtml(ex.notes)}</p>` : ''}
            <table class="sets-table">
                <thead>
                    <tr>
                        <th>Подход</th>
                        <th>Вес (кг)</th>
                        <th>Повторения</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    ${sets.map(set => `
                        <tr>
                            <td>${set.set_number}</td>
                            <td>${set.weight}</td>
                            <td>${set.reps}</td>
                            <td>
                                <button class="btn btn-danger btn-sm delete-set" data-set-id="${set.id}">
                                    🗑️
                                </button>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        </div>
    `;
}

// Отрисовка кардио упражнения
function renderCardioExercise(ex) {
    const cardio = ex.cardio;

    if (!cardio) {
        return `
            <div class="exercise-card">
                <div class="exercise-header">
                    <div>
                        <strong>${escapeHtml(ex.exercise.name)}</strong>
                        <span class="exercise-type cardio">Кардио</span>
                    </div>
                </div>
                <p style="color: var(--text-muted);">Данные не указаны</p>
            </div>
        `;
    }

    return `
        <div class="exercise-card">
            <div class="exercise-header">
                <div>
                    <strong>${escapeHtml(ex.exercise.name)}</strong>
                    <span class="exercise-type cardio">Кардио</span>
                </div>
                <button class="btn btn-danger btn-sm delete-cardio" data-cardio-id="${cardio.id}">
                    🗑️ Удалить
                </button>
            </div>
            ${ex.notes ? `<p style="color: var(--text-muted); margin-bottom: 10px;">${escapeHtml(ex.notes)}</p>` : ''}
            <div class="cardio-data">
                ${cardio.distance_km ? `
                    <div class="cardio-stat">
                        <div class="value">${cardio.distance_km}</div>
                        <div class="label">км</div>
                    </div>
                ` : ''}
                <div class="cardio-stat">
                    <div class="value">${Math.round(cardio.duration_seconds / 60)}</div>
                    <div class="label">минут</div>
                </div>
                ${cardio.avg_heart_rate ? `
                    <div class="cardio-stat">
                        <div class="value">${cardio.avg_heart_rate}</div>
                        <div class="label">уд/мин</div>
                    </div>
                ` : ''}
                ${cardio.calories_burned ? `
                    <div class="cardio-stat">
                        <div class="value">${cardio.calories_burned}</div>
                        <div class="label">ккал</div>
                    </div>
                ` : ''}
            </div>
        </div>
    `;
}

// Обработчики удаления
function attachDeleteHandlers() {
    // Удаление подхода
    document.querySelectorAll('.delete-set').forEach(btn => {
        btn.addEventListener('click', async (e) => {
            e.stopPropagation();
            const setId = btn.getAttribute('data-set-id');

            if (confirm('Удалить этот подход?')) {
                try {
                    await WorkoutAPI.deleteExerciseSet(setId);
                    loadWorkoutDetail(workoutId);
                } catch (error) {
                    alert('Ошибка: ' + error.message);
                }
            }
        });
    });

    // Удаление кардио
    document.querySelectorAll('.delete-cardio').forEach(btn => {
        btn.addEventListener('click', async (e) => {
            e.stopPropagation();
            const cardioId = btn.getAttribute('data-cardio-id');

            if (confirm('Удалить эту кардио запись?')) {
                try {
                    await WorkoutAPI.deleteCardioRecord(cardioId);
                    loadWorkoutDetail(workoutId);
                } catch (error) {
                    alert('Ошибка: ' + error.message);
                }
            }
        });
    });
}

// Загрузка личных рекордов
async function loadPersonalRecords() {
    const container = document.getElementById('records-list');

    try {
        const data = await WorkoutAPI.getPersonalRecords();
        const records = data.records || [];

        if (records.length === 0) {
            container.innerHTML = '<p style="text-align: center; color: var(--text-muted);">У вас пока нет рекордов</p>';
            return;
        }

        container.innerHTML = records.map(record => `
            <div class="record-item">
                <div>
                    <div class="exercise">${escapeHtml(record.exercise.name)}</div>
                    <div class="type">${getRecordTypeLabel(record.record_type)}</div>
                </div>
                <div class="value">${formatRecordValue(record.record_type, record.value)}</div>
            </div>
        `).join('');
    } catch (error) {
        container.innerHTML = `<div class="error-message">${error.message}</div>`;
    }
}

// Вспомогательные функции
function getRecordTypeLabel(type) {
    const labels = {
        'max_weight': 'Максимальный вес',
        'max_reps': 'Максимум повторений',
        'max_distance': 'Максимальная дистанция',
        'best_time': 'Лучшее время'
    };
    return labels[type] || type;
}

function formatRecordValue(type, value) {
    if (type === 'max_weight') return `${value} кг`;
    if (type === 'max_reps') return `${value} раз`;
    if (type === 'max_distance') return `${value} км`;
    if (type === 'best_time') return `${Math.round(value / 60)} мин`;
    return value;
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
