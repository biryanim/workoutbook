document.addEventListener('DOMContentLoaded', () => {
    // Блок для dashboard.html
    if (window.location.pathname.includes('dashboard.html')) {
        console.log("Dashboard page JS running");
        if (!checkAuth()) return;

        loadWorkouts();

        const modal = document.getElementById('createWorkoutModal');
        const btn = document.getElementById('createWorkoutBtn');
        const closeBtn = modal.querySelector('.close');
        const cancelBtn = document.getElementById('cancelCreate');
        const editModal = document.getElementById('editWorkoutModal');
        const closeEditBtn = document.getElementById('closeEditModal');
        const cancelEditBtn = document.getElementById('cancelEdit');

        closeEditBtn.onclick = () => {
            editModal.style.display = 'none';
        };

        cancelEditBtn.onclick = () => {
            editModal.style.display = 'none';
        };

        btn.onclick = () => {
            console.log('Кнопка + Новая тренировка нажата');
            modal.style.display = 'flex';
            document.getElementById('workoutDate').valueAsDate = new Date();
        };

        closeBtn.onclick = () => {
            modal.style.display = 'none';
        };

        cancelBtn.onclick = () => {
            modal.style.display = 'none';
        };

        window.onclick = (e) => {
            if (e.target === modal) {
                modal.style.display = 'none';
            }
            if (e.target === editModal) {
                editModal.style.display = 'none';
            }
            const recordsModal = document.getElementById('recordsModal');
            if (recordsModal && e.target === recordsModal) {
                recordsModal.style.display = 'none';
            }
        };

        document.getElementById('editWorkoutForm').onsubmit = async (e) => {
            e.preventDefault();

            const workoutId = document.getElementById('editWorkoutId').value;
            const name = document.getElementById('editWorkoutName').value;
            const description = document.getElementById('editWorkoutDescription').value;
            const date = document.getElementById('editWorkoutDate').value;

            try {
                await WorkoutAPI.updateWorkout(workoutId, name, description, date);
                editModal.style.display = 'none';
                loadWorkouts();
            } catch (error) {
                alert('Ошибка: ' + error.message);
            }
        };

        document.getElementById('createWorkoutForm').onsubmit = async (e) => {
            e.preventDefault();

            const name = document.getElementById('workoutName').value;
            const description = document.getElementById('workoutDescription').value;
            const date = document.getElementById('workoutDate').value;

            try {
                await WorkoutAPI.createWorkout(name, description, date);
                modal.style.display = 'none';
                e.target.reset();
                loadWorkouts();
            } catch (error) {
                alert('Ошибка: ' + error.message);
            }
        };

        const recordsModal = document.getElementById('recordsModal');
        const recordsLink = document.getElementById('recordsLink');
        const recordsCloseBtn = recordsModal.querySelector('.close');

        recordsLink.onclick = (e) => {
            e.preventDefault();
            recordsModal.style.display = 'flex';
            loadPersonalRecords();
        };

        recordsCloseBtn.onclick = () => {
            recordsModal.style.display = 'none';
        };
    }

    // Блок для workout-detail.html
    if (window.location.pathname.includes('workout-detail.html')) {
        if (!checkAuth()) return;

        const workoutId = new URLSearchParams(window.location.search).get('id');

        if (!workoutId) {
            window.location.href = '/dashboard.html';
        } else {
            loadWorkoutDetail(workoutId);
            initAddExerciseModal(workoutId);

            const deleteWorkoutBtn = document.getElementById('deleteWorkoutBtn');
            if (deleteWorkoutBtn) {
                deleteWorkoutBtn.onclick = async () => {
                    if (confirm('Вы уверены, что хотите удалить эту тренировку? Все упражнения будут удалены.')) {
                        try {
                            await WorkoutAPI.deleteWorkout(workoutId);
                            alert('Тренировка удалена');
                            window.location.href = '/dashboard.html';
                        } catch (error) {
                            alert('Ошибка: ' + error.message);
                        }
                    }
                };
            }
        }
    }
});

// Флаг для предотвращения повторной инициализации
let isExerciseModalInitialized = false;

function initAddExerciseModal(workoutId) {
    if (isExerciseModalInitialized) {
        return; // Если уже инициализировано, выходим
    }
    isExerciseModalInitialized = true;

    const modal = document.getElementById('addExerciseModal');
    const addExerciseBtn = document.getElementById('addExerciseBtn');
    const closeBtn = modal.querySelector('.close');
    const cancelBtn = document.getElementById('cancelAddExercise');
    const saveExerciseBtn = document.getElementById('saveExercise');
    const exerciseSelect = document.getElementById('exerciseSelect');
    const strengthForm = document.getElementById('strengthForm');
    const cardioForm = document.getElementById('cardioForm');
    const addSetBtn = document.getElementById('addSetBtn');

    let exercises = [];
    let setCounter = 0;

    // Функция для добавления подхода
    function addSetRow() {
        console.log('addSetRow вызвана, текущий setCounter:', setCounter);
        console.trace(); // покажет откуда вызывается функция

        setCounter++;
        const container = document.getElementById('sets-container');
        const row = document.createElement('div');
        row.className = 'set-row';
        row.innerHTML = `
            <label>Подход ${setCounter}</label>
            <input type="number" class="set-weight" placeholder="Вес" step="0.5" min="0" required>
            <input type="number" class="set-reps" placeholder="Повт" min="1" required>
            <button type="button" class="remove-set">✕</button>
        `;

        row.querySelector('.remove-set').onclick = () => {
            row.remove();
            renumberSets();
        };

        container.appendChild(row);
    }

    // Функция для перенумерации подходов
    function renumberSets() {
        const rows = document.querySelectorAll('#sets-container .set-row');
        setCounter = rows.length;
        rows.forEach((row, index) => {
            row.querySelector('label').textContent = `Подход ${index + 1}`;
        });
    }

    // Функция очистки формы
    function resetForm() {
        setCounter = 0;
        document.getElementById('sets-container').innerHTML = '';
        document.getElementById('exerciseNotes').value = '';
        document.getElementById('cardioDistance').value = '';
        document.getElementById('cardioDuration').value = '';
        document.getElementById('cardioHeartRate').value = '';
        document.getElementById('cardioCalories').value = '';
        strengthForm.style.display = 'none';
        cardioForm.style.display = 'none';
        exerciseSelect.value = '';
    }

    // Открытие модального окна
    addExerciseBtn.onclick = async () => {
        resetForm();

        try {
            const data = await WorkoutAPI.getExercises();
            exercises = data.exercises || [];

            exerciseSelect.innerHTML = '<option value="">Выберите упражнение</option>' +
                exercises.map(ex => `<option value="${ex.id}">${ex.name} (${ex.type === 'strength' ? 'Силовое' : 'Кардио'})</option>`).join('');

            modal.style.display = 'flex';
        } catch (error) {
            alert('Ошибка загрузки упражнений: ' + error.message);
        }
    };

    // Закрытие модального окна
    closeBtn.onclick = () => {
        modal.style.display = 'none';
        resetForm();
    };

    cancelBtn.onclick = () => {
        modal.style.display = 'none';
        resetForm();
    };

    // Выбор типа упражнения
    exerciseSelect.onchange = () => {
        const selectedExerciseId = parseInt(exerciseSelect.value);
        const selectedExercise = exercises.find(ex => ex.id === selectedExerciseId);

        if (!selectedExercise) {
            strengthForm.style.display = 'none';
            cardioForm.style.display = 'none';
            return;
        }

        if (selectedExercise.type === 'strength') {
            strengthForm.style.display = 'block';
            cardioForm.style.display = 'none';

            if (setCounter === 0) {
                addSetRow();
            }
        } else {
            strengthForm.style.display = 'none';
            cardioForm.style.display = 'block';
        }
    };

    // Кнопка добавления подхода
    addSetBtn.onclick = () => {
        addSetRow();
    };

    // Сохранение упражнения
    saveExerciseBtn.onclick = async () => {
        const selectedExerciseId = parseInt(exerciseSelect.value);
        if (!selectedExerciseId) {
            alert('Выберите упражнение');
            return;
        }

        const selectedExercise = exercises.find(ex => ex.id === selectedExerciseId);
        const notes = document.getElementById('exerciseNotes').value;

        try {
            if (selectedExercise.type === 'strength') {
                const setRows = document.querySelectorAll('#sets-container .set-row');
                if (setRows.length === 0) {
                    alert('Добавьте хотя бы один подход');
                    return;
                }

                const sets = Array.from(setRows).map((row, index) => {
                    const weight = parseFloat(row.querySelector('.set-weight').value);
                    const reps = parseInt(row.querySelector('.set-reps').value);

                    return {
                        set_number: index + 1,
                        weight,
                        reps
                    };
                });

                await WorkoutAPI.addExerciseToWorkout(workoutId, selectedExerciseId, 0, notes, sets);
            } else {
                const distance = parseFloat(document.getElementById('cardioDistance').value) || null;
                const duration = parseInt(document.getElementById('cardioDuration').value);
                const heartRate = parseInt(document.getElementById('cardioHeartRate').value) || null;
                const calories = parseInt(document.getElementById('cardioCalories').value) || null;

                if (!duration) {
                    alert('Укажите время тренировки');
                    return;
                }

                const cardioData = {
                    distance_km: distance,
                    duration_seconds: duration * 60,
                    avg_heart_rate: heartRate,
                    calories_burned: calories
                };

                await WorkoutAPI.addCardioToWorkout(workoutId, selectedExerciseId, 0, notes, cardioData);
            }

            modal.style.display = 'none';
            resetForm();
            loadWorkoutDetail(workoutId);
        } catch (error) {
            alert('Ошибка: ' + error.message);
        }
    };
}
