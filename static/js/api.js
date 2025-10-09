const API_BASE_URL = 'http://localhost:8080/api';

// Функция для получения заголовков с авторизацией
function getAuthHeaders() {
    const token = localStorage.getItem('token');
    const csrfToken = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content');

    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
    };

    if (csrfToken) {
        headers['X-CSRF-Token'] = csrfToken;
    }

    return headers;
}

// API для работы с воркаутами
const WorkoutAPI = {
    async getWorkouts() {
        const response = await fetch(`${API_BASE_URL}/workouts`, {
            headers: getAuthHeaders()
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка загрузки тренировок');
        }

        return await response.json();
    },

    async createWorkout(name, description, date) {
        const response = await fetch(`${API_BASE_URL}/workouts`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({
                name,
                description,
                workout_date: date
            })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка создания тренировки');
        }

        return await response.json();
    },

    async updateWorkout(workoutId, name, description, date) {
        const response = await fetch(`${API_BASE_URL}/workouts/${workoutId}`, {
            method: 'PUT',
            headers: getAuthHeaders(),
            body: JSON.stringify({
                name,
                description,
                workout_date: date
            })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка обновления тренировки');
        }

        return await response.json();
    },

    async deleteWorkout(workoutId) {
        const response = await fetch(`${API_BASE_URL}/workouts/${workoutId}`, {
            method: 'DELETE',
            headers: getAuthHeaders()
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка удаления тренировки');
        }

        return await response.json();
    },

    async getWorkout(workoutId) {
        const response = await fetch(`${API_BASE_URL}/workouts/${workoutId}`, {
            headers: getAuthHeaders()
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка загрузки тренировки');
        }

        return await response.json();
    },

    async addExerciseToWorkout(workoutId, exerciseId, order, notes, sets) {
        const response = await fetch(`${API_BASE_URL}/workouts/${workoutId}/exercises`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({
                exercise_id: exerciseId,
                order: order,
                notes: notes,
                sets: sets
            })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка добавления упражнения');
        }

        return await response.json();
    },

    async addCardioToWorkout(workoutId, exerciseId, order, notes, cardioData) {
        const response = await fetch(`${API_BASE_URL}/workouts/${workoutId}/cardio`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({
                exercise_id: exerciseId,
                order: order,
                notes: notes,
                ...cardioData
            })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка добавления кардио');
        }

        return await response.json();
    },

    async getExercises() {
        const response = await fetch(`${API_BASE_URL}/exercises`, {
            headers: getAuthHeaders()
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка загрузки упражнений');
        }

        return await response.json();
    },

    async deleteExerciseSet(setId) {
        const response = await fetch(`${API_BASE_URL}/sets/${setId}`, {
            method: 'DELETE',
            headers: getAuthHeaders()
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка удаления подхода');
        }

        return await response.json();
    },

    async deleteCardioRecord(cardioId) {
        const response = await fetch(`${API_BASE_URL}/cardio/${cardioId}`, {
            method: 'DELETE',
            headers: getAuthHeaders()
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка удаления кардио записи');
        }

        return await response.json();
    },

    async getPersonalRecords() {
        const response = await fetch(`${API_BASE_URL}/records`, {
            headers: getAuthHeaders()
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка загрузки рекордов');
        }

        return await response.json();
    }
};

const AuthAPI = {
    async register(email, password) {
        const csrfToken = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content');

        const headers = {
            'Content-Type': 'application/json'
        };

        if (csrfToken) {
            headers['X-CSRF-Token'] = csrfToken;
        }

        const response = await fetch(`${API_BASE_URL}/register`, {
            method: 'POST',
            headers: headers,
            body: JSON.stringify({ email, password })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка регистрации');
        }

        return await response.json();
    },

    async login(email, password) {
        const csrfToken = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content');

        const headers = {
            'Content-Type': 'application/json'
        };

        if (csrfToken) {
            headers['X-CSRF-Token'] = csrfToken;
        }

        const response = await fetch(`${API_BASE_URL}/login`, {
            method: 'POST',
            headers: headers,
            body: JSON.stringify({ email, password })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка входа');
        }

        const data = await response.json();
        localStorage.setItem('token', data.token);

        return data;
    },

    logout() {
        localStorage.removeItem('token');
        window.location.href = '/login.html';
    }
};

