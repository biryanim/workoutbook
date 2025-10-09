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
        console.log('Модальное окно уже инициализировано');
        return;
    }
    isExerciseModalInitialized = true;
    console.log('Инициализация модального окна добавления упражнения');

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
    let allExercises = [];

    function addSetRow() {
        const container = document.getElementById('sets-container');
        const currentCount = container.querySelectorAll('.set-row').length;
        const newNumber = currentCount + 1;

        const row = document.createElement('div');
        row.className = 'set-row';
        row.innerHTML = `
            <label>Подход ${newNumber}</label>
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

    function renumberSets() {
        const container = document.getElementById('sets-container');
        const rows = container.querySelectorAll('.set-row');

        rows.forEach((row, index) => {
            row.querySelector('label').textContent = `Подход ${index + 1}`;
        });
    }

    function resetForm() {
        const container = document.getElementById('sets-container');
        container.innerHTML = '';
        document.getElementById('exerciseNotes').value = '';

        const distanceInput = document.getElementById('cardioDistance');
        const durationInput = document.getElementById('cardioDuration');
        if (distanceInput) distanceInput.value = '';
        if (durationInput) durationInput.value = '';

        strengthForm.style.display = 'none';
        cardioForm.style.display = 'none';
        exerciseSelect.value = '';
    }

    function populateExercises(exercisesList) {
        console.log('=== populateExercises вызвана ===');
        console.log('Получено упражнений:', exercisesList.length);

        exerciseSelect.innerHTML = '<option value="">Выберите упражнение</option>';

        const strengthExercises = exercisesList.filter(ex => ex.type === 'strength');
        const cardioExercises = exercisesList.filter(ex => ex.type === 'cardio');

        console.log('Силовых:', strengthExercises.length);
        console.log('Кардио:', cardioExercises.length);

        if (strengthExercises.length > 0) {
            const strengthGroup = document.createElement('optgroup');
            strengthGroup.label = '💪 Силовые упражнения';

            const byMuscle = strengthExercises.reduce((acc, ex) => {
                const group = ex.muscle_group || 'Другое';
                if (!acc[group]) acc[group] = [];
                acc[group].push(ex);
                return acc;
            }, {});

            Object.keys(byMuscle).sort().forEach(muscle => {
                byMuscle[muscle].forEach(ex => {
                    const option = document.createElement('option');
                    option.value = ex.id;
                    option.textContent = `${ex.name} (${muscle})`;
                    strengthGroup.appendChild(option);
                });
            });

            exerciseSelect.appendChild(strengthGroup);
        }

        if (cardioExercises.length > 0) {
            const cardioGroup = document.createElement('optgroup');
            cardioGroup.label = '🏃 Кардио упражнения';

            cardioExercises.forEach(ex => {
                const option = document.createElement('option');
                option.value = ex.id;
                option.textContent = ex.name;
                cardioGroup.appendChild(option);
            });

            exerciseSelect.appendChild(cardioGroup);
        }

        console.log('Итого options в select:', exerciseSelect.options.length);
    }

    function filterExercises(filterType) {
        console.log('Фильтрация по типу:', filterType);

        if (filterType === 'all') {
            exercises = allExercises;
        } else {
            exercises = allExercises.filter(ex => ex.type === filterType);
        }

        console.log('Упражнений после фильтрации:', exercises.length);
        populateExercises(exercises);
    }

    // Обработчик фильтра
    const filterRadios = document.querySelectorAll('input[name="exerciseFilter"]');
    if (filterRadios.length > 0) {
        filterRadios.forEach(radio => {
            radio.onchange = () => {
                // ВАЖНО: Проверяем, что данные загружены
                if (allExercises.length === 0) {
                    console.log('Данные ещё не загружены, пропускаем фильтрацию');
                    return;
                }

                filterExercises(radio.value);
                exerciseSelect.value = '';
                strengthForm.style.display = 'none';
                cardioForm.style.display = 'none';
            };
        });
    }

    addExerciseBtn.onclick = async () => {
        console.log('=== Кнопка "Добавить упражнение" нажата ===');
        resetForm();

        try {
            console.log('Запрос к API...');
            const data = await WorkoutAPI.getExercises();
            console.log('Данные получены:', data);

            allExercises = data.exercises || [];
            exercises = allExercises;
            console.log('Всего упражнений:', allExercises.length);

            populateExercises(exercises);

            console.log('После populateExercises, options:', exerciseSelect.options.length);

            // Устанавливаем фильтр "Все" ПОСЛЕ загрузки
            const allFilter = document.querySelector('input[name="exerciseFilter"][value="all"]');
            if (allFilter) {
                allFilter.checked = true;
            }

            modal.style.display = 'flex';
        } catch (error) {
            console.error('Ошибка:', error);
            alert('Ошибка загрузки упражнений: ' + error.message);
        }
    };

    closeBtn.onclick = () => {
        modal.style.display = 'none';
        resetForm();
    };

    cancelBtn.onclick = () => {
        modal.style.display = 'none';
        resetForm();
    };

    exerciseSelect.onchange = () => {
        const selectedExerciseId = parseInt(exerciseSelect.value);
        const selectedExercise = exercises.find(ex => ex.id === selectedExerciseId);

        console.log('Выбрано упражнение:', selectedExercise);

        if (!selectedExercise) {
            strengthForm.style.display = 'none';
            cardioForm.style.display = 'none';
            return;
        }

        if (selectedExercise.type === 'strength') {
            strengthForm.style.display = 'block';
            cardioForm.style.display = 'none';

            const container = document.getElementById('sets-container');
            if (container.querySelectorAll('.set-row').length === 0) {
                addSetRow();
            }
        } else {
            strengthForm.style.display = 'none';
            cardioForm.style.display = 'block';
        }
    };

    addSetBtn.onclick = () => {
        addSetRow();
    };

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
                const container = document.getElementById('sets-container');
                const setRows = container.querySelectorAll('.set-row');

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

                if (!duration) {
                    alert('Укажите время тренировки');
                    return;
                }

                const cardioData = {
                    distance_km: distance,
                    duration_seconds: duration * 60,
                    avg_heart_rate: null,
                    calories_burned: null
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

