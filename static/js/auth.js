// Проверка, авторизован ли пользователь
function isAuthenticated() {
    const token = localStorage.getItem('token');
    return !!token; // вернёт true, если токен есть, иначе false
}

// Проверка авторизации на защищенных страницах
function checkAuth() {
    if (!isAuthenticated()) {
        window.location.href = '/login.html';
        return false;
    }
    return true;
}

// Обработка формы регистрации
const registerForm = document.getElementById('registerForm');
if (registerForm) {
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const confirmPassword = document.getElementById('confirmPassword').value;
        const errorDiv = document.getElementById('error-message');

        // Проверка совпадения паролей
        if (password !== confirmPassword) {
            errorDiv.textContent = 'Пароли не совпадают';
            errorDiv.style.display = 'block';
            return;
        }

        try {
            await AuthAPI.register(email, password);
            alert('Регистрация успешна! Теперь войдите в систему.');
            window.location.href = '/login.html';
        } catch (error) {
            errorDiv.textContent = error.message;
            errorDiv.style.display = 'block';
        }
    });
}

// Обработка формы входа
const loginForm = document.getElementById('loginForm');
if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        try {
            await AuthAPI.login(email, password);
            window.location.href = '/dashboard.html';
        } catch (error) {
            const errorDiv = document.getElementById('error-message');
            errorDiv.textContent = error.message;
            errorDiv.style.display = 'block';
        }
    });
}

// Кнопка выхода
document.addEventListener('DOMContentLoaded', () => {
    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', (e) => {
            e.preventDefault();
            if (confirm('Вы уверены, что хотите выйти?')) {
                AuthAPI.logout();
            }
        });
    }
});
