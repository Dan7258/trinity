// static/js/auth.js

// === Загрузка хедера ===
document.addEventListener('DOMContentLoaded', () => {
    const container = document.getElementById('header-container');
    if (!container) return;

    fetch('/static/html/header.html')
        .then(res => {
            if (!res.ok) throw new Error(`HTTP ${res.status}`);
            return res.text();
        })
        .then(html => {
            container.innerHTML = html;
            updateAuthDisplay();
        })
        .catch(err => {
            container.innerHTML = `<p style="color:red;">Ошибка: ${err.message}</p>`;
        });
});

// === Обновление хедера ===
function updateAuthDisplay() {
    const token = localStorage.getItem('token');
    const authButtons = document.getElementById('auth-buttons');
    const userInfo = document.getElementById('user-info');

    if (token && authButtons && userInfo) {
        authButtons.style.display = 'none';
        userInfo.style.display = 'flex';
    } else if (authButtons && userInfo) {
        authButtons.style.display = 'flex';
        userInfo.style.display = 'none';
    }
    updateHistoryTab();
}

// === Кнопки ===
function goToHome() { window.location.href = '/index.html'; }
function goToLogin() { window.location.href = '/login.html'; }
function goToRegister() { window.location.href = '/register.html'; }

async function goToAdmin() {
    const token = localStorage.getItem('token');
    if (!token) {
        alert('Требуется авторизация');
        window.location.href = '/login.html';
        return;
    }

    try {
        // Запрашиваем страницу с токеном в заголовке
        const response = await fetch('/admin/history/', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (response.ok) {
            // Получаем HTML
            const html = await response.text();

            // Заменяем текущую страницу на полученный HTML
            document.open();
            document.write(html);
            document.close();

            // Обновляем URL без перезагрузки
            window.history.pushState({}, '', '/admin/history/');

        } else if (response.status === 401) {
            alert('Сессия истекла. Войдите заново.');
            localStorage.removeItem('token');
            window.location.href = '/login.html';
        } else if (response.status === 403) {
            alert('Доступ запрещён. Недостаточно прав.');
            window.location.href = '/index.html';
        } else {
            alert('Ошибка загрузки страницы');
        }
    } catch (err) {
        console.error('Ошибка:', err);
        alert('Сетевая ошибка');
    }
}

function logout() {
    localStorage.removeItem('token');
    updateAuthDisplay();
    window.location.href = '/';
}

// === Синхронизация между вкладками ===
window.addEventListener('storage', (e) => {
    if (e.key === 'token') {
        updateAuthDisplay();
    }
});

// === Показ/скрытие вкладки истории ===
function updateHistoryTab() {
    const historyTab = document.getElementById('history-tab');
    const token = localStorage.getItem('token');
    if (historyTab) {
        historyTab.style.display = token ? 'block' : 'none';
    }
}

document.addEventListener('DOMContentLoaded', updateHistoryTab);
window.addEventListener('storage', (e) => {
    if (e.key === 'token') {
        updateAuthDisplay();
        updateHistoryTab();
    }
});