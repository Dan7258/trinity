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
function goToLogin() { window.location.href = '/login.html'; }
function goToRegister() { window.location.href = '/register.html'; }

// НОВАЯ ФУНКЦИЯ
function goToAdmin() {
    const token = localStorage.getItem('token');
    if (token) {
        // Передаём токен как query-параметр — сервер должен его проверить
        window.location.href = `/admin.html?token=${encodeURIComponent(token)}`;
    } else {
        alert('Требуется авторизация');
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
        updateHistoryTab();
    }
});