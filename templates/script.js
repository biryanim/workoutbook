// Глобальные переменные
let authToken = localStorage.getItem('auth_token');
let currentUser = JSON.parse(localStorage.getItem('current_user') || '{}');
let currentWorkoutId = null;
let exercises = [];

// API базовый URL
const API_BASE = '/api';

// Утилитные функции
function showAlert(message, type = 'success') {
    const alertDiv = document.createElement('div');
    alertDiv.className = `alert alert-${type}`;
    alertDiv.textContent = message;

    const container = document.querySelector('.container') || document.body;
    container.insertBefore(alertDiv, container.firstChild);

    setTimeout(() => alertDiv.remove(), 5000);
}

function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}

async function apiRequest(endpoint, options = {}) {
    const url = `${API_BASE}${endpoint}`;
    const config = {
        headers: {
            'Content-Type': 'application/json',
            ...(authToken && { 'Authorization': `Bearer ${authToken}` })
        },
        ...options
    };

    try {
        const response = await fetch(url, config);
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Ошибка запроса');
        }

        return data;
    } catch (error) {
        console.error('API Error:', error);
        throw error;
    }
}

// Аутентификация
function initAuth() {
    if (authToken && currentUser.id) {
        showApp();
    } else {
        showAuth();
    }
}

function showAuth() {
    document.getElementById('auth-container').classList.remove('hidden');
    document.getElementById('app-container').classList.add('hidden');
}

function showApp() {
    document.getElementById('auth-container').classList.add('hidden');
    document.getElementById('app-container').classList.remove('hidden');
    document.getElementById('user-name').textContent = currentUser.username;
    loadData();
}

// Обработчики форм
document.getElementById('auth-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const isLogin = document.getElementById('auth-submit').textContent === 'Войти';
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const email = document.getElementById('email').value;

    try {
        if (isLogin) {
            const result = await apiRequest('/login', {
                method: 'POST',
                body: JSON.stringify({ email, password })
            });

            authToken = result.token;
            currentUser = result.user;
            localStorage.setItem('auth_token', authToken);
            localStorage.setItem('current_user', JSON.stringify(currentUser));
            showApp();
        } else {
            await apiRequest('/register', {
                method: 'POST',
                body: JSON.stringify({ username, email, password })
            });
            showAlert('Регистрация успешна! Теперь войдите в систему.');
            toggleAuthMode();
        }
    } catch (error) {
        showAlert(error.message, 'error');
    }
});

document.getElementById('auth-toggle').addEventListener('click', toggleAuthMode);

function toggleAuthMode() {
    const isLogin = document.getElementById('auth-submit').textContent === 'Войти';

    if (isLogin) {
        document.getElementById('auth-title').textContent = 'Регистрация';
        document.getElementById('auth-submit').textContent = 'Зарегистрироваться';
        document.getElementById('auth-toggle').textContent = 'Вход';
        document.getElementById('email-group').classList.remove('hidden');
    } else {
        document.getElementById('auth-title').textContent = 'Вход в систему';
        document.getElementById('auth-submit').textContent = 'Войти';
        document.getElementById('auth-toggle').textContent = 'Регистрация';
        document.getElementById('email-group').classList.add('hidden');
    }
}

// Выход из системы
document.getElementById('logout-btn').addEventListener('click', () => {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('current_user');
    authToken = null;
    currentUser = {};
    showAuth();
});

// Управление вкладками
document.querySelectorAll('.tab').forEach(tab => {
    tab.addEventListener('click', () => {
        const tabName = tab.dataset.tab;

        // Переключение активной вкладки
        document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(tc => tc.classList.remove('active'));

        tab.classList.add('active');
        document.getElementById(`${tabName}-tab`).classList.add('active');

        // Загрузка данных для вкладки
        if (tabName === 'exercises') {
            loadExercises();
        } else if (tabName === 'records') {
            loadPersonalRecords();
        } else if (tabName === 'progress') {
            loadProgressExercises();
        }
    });
});

// Загрузка данных
async function loadData() {
    await loadWorkouts();
    await loadExercises();
}

async function loadWorkouts() {
    try {
        const workouts = await apiRequest('/workouts');
        displayWorkouts(workouts);
    } catch (error) {
        document.getElementById('workouts-list').innerHTML = `<p>Ошибка загрузки: ${error.message}</p>`;
    }
}

function displayWorkouts(workouts) {
    const container = document.getElementById('workouts-list');

    if (workouts.length === 0) {
        container.innerHTML = '<p>У вас пока нет тренировок. Создайте первую!</p>';
        return;
    }

    container.innerHTML = workouts.map(workout => `
                <div class="workout-item">
                    <div>
                        <h4>${workout.name}</h4>
                        <p>${formatDate(workout.date)}</p>
                        ${workout.notes ? `<p><em>${workout.notes}</em></p>` : ''}
                    </div>
                    <div>
                        <button class="btn btn-success btn-small" onclick="viewWorkout(${workout.id})">Подробнее</button>
                    </div>
                </div>
            `).join('');
}

async function loadExercises() {
    try {
        const type = document.getElementById('exercise-type-filter')?.value || '';
        const url = type ? `/exercises?type=${type}` : '/exercises';
        exercises = await apiRequest(url);
        displayExercises(exercises);
        updateExerciseSelects();
    } catch (error) {
        document.getElementById('exercises-list').innerHTML = `<p>Ошибка загрузки: ${error.message}</p>`;
    }
}

function displayExercises(exercisesList) {
    const container = document.getElementById('exercises-list');

    if (exercisesList.length === 0) {
        container.innerHTML = '<p>Упражнения не найдены.</p>';
        return;
    }

    const grouped = exercisesList.reduce((acc, exercise) => {
        const group = exercise.muscle_group || 'Другое';
        if (!acc[group]) acc[group] = [];
        acc[group].push(exercise);
        return acc;
    }, {});

    container.innerHTML = Object.entries(grouped).map(([group, exs]) => `
                <div style="margin-bottom: 25px;">
                    <h4 style="color: #667eea; margin-bottom: 15px;">${group}</h4>
                    ${exs.map(exercise => `
                        <div class="exercise-item">
                            <div style="display: flex; justify-content: space-between; align-items: center;">
                                <div style="flex: 1;">
                                    <h5>${exercise.name}</h5>
                                    <p style="margin: 5px 0; color: #666;">${exercise.description}</p>
                                    <span style="background: ${exercise.type === 'strength' ? '#28a745' : '#17a2b8'}; color: white; padding: 3px 8px; border-radius: 12px; font-size: 12px;">
                                        ${exercise.type === 'strength' ? 'Силовое' : 'Кардио'}
                                    </span>
                                </div>
                            </div>
                        </div>
                    `).join('')}
                </div>
            `).join('');
}

function updateExerciseSelects() {
    const selects = ['exercise-select', 'progress-exercise'];

    selects.forEach(selectId => {
        const select = document.getElementById(selectId);
        if (select) {
            const currentValue = select.value;
            select.innerHTML = '<option value="">Выберите упражнение</option>' +
                exercises.map(ex => `<option value="${ex.id}" data-type="${ex.type}">${ex.name}</option>`).join('');
            select.value = currentValue;
        }
    });
}

// Обработчик фильтра упражнений
document.getElementById('exercise-type-filter').addEventListener('change', loadExercises);

// Создание тренировки
document.getElementById('create-workout-btn').addEventListener('click', () => {
    document.getElementById('workout-date').value = new Date().toISOString().slice(0, 16);
    document.getElementById('workout-modal').classList.add('show');
});

document.getElementById('close-workout-modal').addEventListener('click', () => {
    document.getElementById('workout-modal').classList.remove('show');
});

document.getElementById('workout-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const workoutData = {
        name: document.getElementById('workout-name').value,
        date: new Date(document.getElementById('workout-date').value).toISOString(),
        notes: document.getElementById('workout-notes').value
    };

    try {
        await apiRequest('/workouts', {
            method: 'POST',
            body: JSON.stringify(workoutData)
        });

        showAlert('Тренировка создана успешно!');
        document.getElementById('workout-modal').classList.remove('show');
        document.getElementById('workout-form').reset();
        loadWorkouts();
    } catch (error) {
        showAlert(error.message, 'error');
    }
});

// Просмотр тренировки
async function viewWorkout(workoutId) {
    currentWorkoutId = workoutId;

    try {
        const data = await apiRequest(`/workouts/${workoutId}`);
        displayWorkoutDetails(data.workout, data.exercises);
    } catch (error) {
        showAlert(error.message, 'error');
    }
}

function displayWorkoutDetails(workout, workoutExercises) {
    const container = document.getElementById('workouts-list');

    container.innerHTML = `
                <div class="card">
                    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
                        <div>
                            <h3>${workout.name}</h3>
                            <p>${formatDate(workout.date)}</p>
                            ${workout.notes ? `<p><em>${workout.notes}</em></p>` : ''}
                        </div>
                        <div>
                            <button class="btn btn-primary btn-small" onclick="showAddExerciseModal()">Добавить упражнение</button>
                            <button class="btn btn-secondary btn-small" onclick="loadWorkouts()">Назад</button>
                        </div>
                    </div>

                    <h4>Упражнения:</h4>
                    <div id="workout-exercises">
                        ${workoutExercises.length === 0 ?
        '<p>В этой тренировке пока нет упражнений.</p>' :
        workoutExercises.map(we => `
                                <div class="exercise-item">
                                    <h5>${we.exercise.name}</h5>
                                    <div class="exercise-stats">
                                        ${we.exercise.type === 'strength' ? `
                                            <div class="stat-item">Подходы: ${we.sets}</div>
                                            <div class="stat-item">Повторения: ${we.reps}</div>
                                            <div class="stat-item">Вес: ${we.weight} кг</div>
                                        ` : `
                                            <div class="stat-item">Время: ${Math.floor(we.duration / 60)} мин</div>
                                            ${we.distance > 0 ? `<div class="stat-item">Расстояние: ${we.distance} км</div>` : ''}
                                        `}
                                    </div>
                                </div>
                            `).join('')
    }
                    </div>
                </div>
            `;
}

// Добавление упражнения в тренировку
function showAddExerciseModal() {
    updateExerciseSelects();
    document.getElementById('exercise-modal').classList.add('show');
}

document.getElementById('close-exercise-modal').addEventListener('click', () => {
    document.getElementById('exercise-modal').classList.remove('show');
});

// Переключение полей в зависимости от типа упражнения
document.getElementById('exercise-select').addEventListener('change', (e) => {
    const selectedOption = e.target.options[e.target.selectedIndex];
    const exerciseType = selectedOption.dataset.type;

    if (exerciseType === 'cardio') {
        document.getElementById('strength-fields').classList.add('hidden');
        document.getElementById('cardio-fields').classList.remove('hidden');
    } else {
        document.getElementById('strength-fields').classList.remove('hidden');
        document.getElementById('cardio-fields').classList.add('hidden');
    }
});

document.getElementById('exercise-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const selectedOption = document.getElementById('exercise-select').options[document.getElementById('exercise-select').selectedIndex];
    const exerciseType = selectedOption.dataset.type;

    const exerciseData = {
        exercise_id: parseInt(document.getElementById('exercise-select').value),
        sets: parseInt(document.getElementById('sets').value)
    };

    if (exerciseType === 'strength') {
        exerciseData.reps = parseInt(document.getElementById('reps').value);
        exerciseData.weight = parseFloat(document.getElementById('weight').value);
        exerciseData.duration = 0;
        exerciseData.distance = 0;
    } else {
        exerciseData.reps = 0;
        exerciseData.weight = 0;
        exerciseData.duration = parseInt(document.getElementById('duration').value) * 60; // в секундах
        exerciseData.distance = parseFloat(document.getElementById('distance').value);
    }

    try {
        await apiRequest(`/workouts/${currentWorkoutId}/exercises`, {
            method: 'POST',
            body: JSON.stringify(exerciseData)
        });

        showAlert('Упражнение добавлено!');
        document.getElementById('exercise-modal').classList.remove('show');
        document.getElementById('exercise-form').reset();
        viewWorkout(currentWorkoutId); // Обновить детали тренировки
    } catch (error) {
        showAlert(error.message, 'error');
    }
});

// Личные рекорды
async function loadPersonalRecords() {
    try {
        const records = await apiRequest('/records');
        displayPersonalRecords(records);
    } catch (error) {
        document.getElementById('records-list').innerHTML = `<p>Ошибка загрузки: ${error.message}</p>`;
    }
}

function displayPersonalRecords(records) {
    const container = document.getElementById('records-list');

    if (records.length === 0) {
        container.innerHTML = '<p>У вас пока нет личных рекордов. Добавьте упражнения в тренировки!</p>';
        return;
    }

    container.innerHTML = `
                <div class="records-grid">
                    ${records.map(record => `
                        <div class="record-card">
                            <h4>${record.exercise.name}</h4>
                            <div style="font-size: 2em; margin: 10px 0;">
                                ${record.weight}кг × ${record.reps}
                            </div>
                            <p>Установлен: ${formatDate(record.date)}</p>
                            <small>Расчетный 1RM: ${(record.weight * (1 + record.reps / 30)).toFixed(1)}кг</small>
                        </div>
                    `).join('')}
                </div>
            `;
}

// Прогресс
async function loadProgressExercises() {
    updateExerciseSelects();
}

document.getElementById('progress-exercise').addEventListener('change', async (e) => {
    const exerciseId = e.target.value;
    if (!exerciseId) {
        document.getElementById('progress-chart').innerHTML = '<p>Выберите упражнение для просмотра прогресса</p>';
        return;
    }

    try {
        const progress = await apiRequest(`/progress?exercise_id=${exerciseId}`);
        displayProgress(progress);
    } catch (error) {
        document.getElementById('progress-chart').innerHTML = `<p>Ошибка загрузки прогресса: ${error.message}</p>`;
    }
});

function displayProgress(progressData) {
    const container = document.getElementById('progress-chart');

    if (progressData.length === 0) {
        container.innerHTML = '<p>Нет данных для отображения прогресса по этому упражнению</p>';
        return;
    }

    // Простая визуализация прогресса
    const maxWeight = Math.max(...progressData.map(p => p.weight));
    const chartHtml = progressData.map((point, index) => {
        const height = (point.weight / maxWeight) * 200;
        const date = new Date(point.date).toLocaleDateString('ru-RU');

        return `
                    <div style="display: inline-block; margin: 0 5px; text-align: center;">
                        <div style="height: 200px; display: flex; align-items: end;">
                            <div style="
                                background: linear-gradient(to top, #667eea, #764ba2);
                                width: 30px;
                                height: ${height}px;
                                border-radius: 3px 3px 0 0;
                            "></div>
                        </div>
                        <div style="font-size: 12px; margin-top: 5px;">
                            <div>${point.weight}кг</div>
                            <div>${date}</div>
                        </div>
                    </div>
                `;
    }).join('');

    container.innerHTML = `
                <div style="padding: 20px; text-align: center;">
                    <h4 style="margin-bottom: 20px;">Прогресс по весу</h4>
                    <div style="display: flex; justify-content: center; align-items: end; overflow-x: auto;">
                        ${chartHtml}
                    </div>
                </div>
            `;
}

// Закрытие модальных окон по клику вне их
document.querySelectorAll('.modal').forEach(modal => {
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.classList.remove('show');
        }
    });
});

// Инициализация приложения
document.addEventListener('DOMContentLoaded', () => {
    initAuth();
});