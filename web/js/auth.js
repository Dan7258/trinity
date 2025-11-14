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
    const token = localStorage.getItem('token'); // ← token
    const authButtons = document.getElementById('auth-buttons');
    const userInfo = document.getElementById('user-info');

    if (token && authButtons && userInfo) {
        authButtons.style.display = 'none';
        userInfo.style.display = 'block';
    } else if (authButtons && userInfo) {
        authButtons.style.display = 'flex';
        userInfo.style.display = 'none';
    }
    updateHistoryTab();
}

// === Кнопки ===
function goToLogin() { window.location.href = '/login.html'; }
function goToRegister() { window.location.href = '/register.html'; }

function logout() {
    localStorage.removeItem('token'); // ← token
    updateAuthDisplay();
    window.location.href = '/';
}

// === Синхронизация между вкладками ===
window.addEventListener('storage', (e) => {
    if (e.key === 'token') { // ← token
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

// Вызываем при загрузке и обновлении
document.addEventListener('DOMContentLoaded', updateHistoryTab);
window.addEventListener('storage', (e) => {
    if (e.key === 'token') {
        updateHistoryTab();
    }
});