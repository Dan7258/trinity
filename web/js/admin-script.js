// static/js/admin-script.js

let currentToken = null;
let currentLogin = '';
let userIdToDelete = null;

// Извлекаем токен из localStorage
function getToken() {
    return localStorage.getItem('token');
}

async function loadAdminHistory() {
    const login = document.getElementById('admin-login-input').value.trim();
    if (!login) {
        alert('Введите логин пользователя');
        return;
    }

    const token = getToken();
    if (!token) {
        document.getElementById('admin-status').innerHTML =
            '<span style="color:red;">Сессия истекла. <a href="/login.html">Войдите заново</a></span>';
        return;
    }

    currentToken = token;
    currentLogin = login;

    const statusEl = document.getElementById('admin-status');
    statusEl.textContent = 'Загрузка...';

    const activeTab = document.querySelector('.tab-button.active').dataset.tab;
    await fetchAdminHistory(activeTab, login);
}

// Универсальная загрузка истории для админа
async function fetchAdminHistory(algorithm, login) {
    const contentEl = document.getElementById(`${algorithm}-history-content`);
    contentEl.innerHTML = '<p><em>Загрузка...</em></p>';

    try {
        const response = await fetch(`/admin/history/${algorithm}/${encodeURIComponent(login)}`, {
            headers: {
                'Authorization': `Bearer ${currentToken}`,
                'Content-Type': 'application/json'
            }
        });

        if (!response.ok) {
            if (response.status === 403) {
                contentEl.innerHTML = '<p style="color:red;">Доступ запрещён (недостаточно прав)</p>';
            } else if (response.status === 404) {
                contentEl.innerHTML = '<p style="color:#666;">Пользователь не найден или история пуста</p>';
            } else {
                contentEl.innerHTML = `<p style="color:red;">Ошибка: ${response.status}</p>`;
            }
            return;
        }

        const data = await response.json();
        renderAdminHistory(data, algorithm);

    } catch (err) {
        contentEl.innerHTML = `<p style="color:red;">Ошибка сети: ${err.message}</p>`;
    }
}

// Рендер истории
function renderAdminHistory(data, algorithm) {
    const contentEl = document.getElementById(`${algorithm}-history-content`);
    if (!Array.isArray(data) || data.length === 0) {
        contentEl.innerHTML = '<p style="color:#666;">История пуста</p>';
        return;
    }

    data.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));

    let html = `<table class="history-table"><thead><tr><th>Время</th><th>Результат</th><th>Действие</th></tr></thead><tbody>`;

    data.forEach(item => {
        const time = new Date(item.created_at).toLocaleString('ru-RU');
        let preview = '', fullText = '', decryptPayload = null;

        if (algorithm === 'kuznechik') {
            const enc = item.encrypted_message || '';
            const key = item.key || '';
            fullText = `Шифротекст: ${enc}\nКлюч: ${key}`;
            preview = enc.substring(0, 60) + (enc.length > 60 ? '...' : '');
            if (enc && key) decryptPayload = { encrypted_message: enc, key };
        } else if (algorithm === 'rsa') {
            const enc = item.encrypted_message || '';
            const d = item.d || '';
            const n = item.n || '';
            fullText = `Зашифровано: ${enc}\nd: ${d}\nn: ${n}`;
            preview = enc.substring(0, 60) + (enc.length > 60 ? '...' : '');
            if (enc && d && n) decryptPayload = { encrypted_message: enc, d, n };
        } else if (algorithm === 'stribog') {
            fullText = item.encrypted_message || '';
            preview = fullText.substring(0, 60) + (fullText.length > 60 ? '...' : '');
        }

        html += `
            <tr>
                <td style="white-space: nowrap;">${time}</td>
                <td class="full-text-cell" title="Клик — показать, двойной клик — копировать">
                    <div class="text-preview">${preview || '—'}</div>
                    <div class="text-full" style="display:none; margin-top:8px; padding:8px; background:#f8f9fa; border-radius:4px; white-space:pre-wrap; word-break:break-all; font-family:monospace; font-size:12px;">
                        ${fullText}
                    </div>
                </td>
                <td>
                    ${decryptPayload ? `
                        <button class="copy-btn" style="font-size:11px;padding:4px 8px;" 
                                onclick='decryptFromHistory("${algorithm}", ${JSON.stringify(decryptPayload).replace(/'/g, "&#39;")})'>
                            Расшифровать
                        </button>
                    ` : '—'}
                </td>
            </tr>`;
    });

    html += `</tbody></table>`;
    contentEl.innerHTML = html;

    document.querySelectorAll('.full-text-cell').forEach(cell => {
        cell.addEventListener('click', function() {
            const full = this.querySelector('.text-full');
            const prev = this.querySelector('.text-preview');
            full.style.display = full.style.display === 'none' ? 'block' : 'none';
            prev.style.display = prev.style.display === 'none' ? 'block' : 'none';
        });
        cell.addEventListener('dblclick', function(e) {
            e.stopPropagation();
            navigator.clipboard.writeText(this.querySelector('.text-full').innerText);
            alert('Скопировано!');
        });
    });
}

// === УДАЛЕНИЕ ПОЛЬЗОВАТЕЛЯ ===

// Показать модальное окно подтверждения
function showAdminDeleteConfirmation() {
    const userIdInput = document.getElementById('admin-delete-user-id');
    const userId = userIdInput.value.trim();

    if (!userId || isNaN(userId) || parseInt(userId) <= 0) {
        alert('Введите корректный ID пользователя');
        return;
    }

    const token = getToken();
    if (!token) {
        alert('Сессия истекла. Войдите заново.');
        window.location.href = '/login.html';
        return;
    }

    userIdToDelete = parseInt(userId);
    document.getElementById('confirm-user-id').textContent = userIdToDelete;
    document.getElementById('admin-delete-confirmation-modal').style.display = 'flex';
}

// Закрыть модальное окно
function closeAdminDeleteConfirmation() {
    document.getElementById('admin-delete-confirmation-modal').style.display = 'none';
    userIdToDelete = null;
}

// Выполнить удаление
async function executeAdminDeleteUser() {
    if (!userIdToDelete) {
        alert('ID пользователя не указан');
        closeAdminDeleteConfirmation();
        return;
    }

    const token = getToken();
    if (!token) {
        alert('Сессия истекла');
        window.location.href = '/login.html';
        return;
    }

    const statusEl = document.getElementById('admin-delete-status');
    statusEl.textContent = 'Удаление...';
    statusEl.style.color = '#666';

    try {
        const response = await fetch(`/api/delete-user/${userIdToDelete}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        });

        if (response.ok) {
            statusEl.textContent = `Пользователь с ID ${userIdToDelete} успешно удалён`;
            statusEl.style.color = 'green';

            // Очищаем поле ввода
            document.getElementById('admin-delete-user-id').value = '';

            // Закрываем модалку
            closeAdminDeleteConfirmation();

            // Очищаем статус через 3 секунды
            setTimeout(() => {
                statusEl.textContent = '';
            }, 3000);

        } else if (response.status === 401) {
            statusEl.textContent = 'Сессия истекла';
            statusEl.style.color = 'red';
            setTimeout(() => {
                window.location.href = '/login.html';
            }, 1500);

        } else if (response.status === 403) {
            statusEl.textContent = 'Доступ запрещён (недостаточно прав)';
            statusEl.style.color = 'red';
            closeAdminDeleteConfirmation();

        } else if (response.status === 404) {
            statusEl.textContent = 'Пользователь не найден';
            statusEl.style.color = 'red';
            closeAdminDeleteConfirmation();

        } else {
            const data = await response.json().catch(() => ({}));
            statusEl.textContent = `Ошибка: ${data.message || 'Не удалось удалить пользователя'}`;
            statusEl.style.color = 'red';
            closeAdminDeleteConfirmation();
        }

    } catch (err) {
        console.error('Ошибка удаления:', err);
        statusEl.textContent = `Сетевая ошибка: ${err.message}`;
        statusEl.style.color = 'red';
        closeAdminDeleteConfirmation();
    }
}

// Переключение вкладок
document.querySelectorAll('.tab-button').forEach(btn => {
    btn.addEventListener('click', () => {
        document.querySelectorAll('.tab-button').forEach(b => b.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
        btn.classList.add('active');
        document.getElementById(btn.dataset.tab).classList.add('active');

        if (currentLogin) {
            fetchAdminHistory(btn.dataset.tab, currentLogin);
        }
    });
});

// Функции модалки расшифровки
function decryptFromHistory(algorithm, payload) {
    // Используем функцию из script.js
    if (typeof window.decryptFromHistory === 'function') {
        window.decryptFromHistory(algorithm, payload);
    }
}

function closeDecryptModal() {
    document.getElementById('decrypt-modal').style.display = 'none';
}

async function copyModalResult() {
    const text = document.getElementById('decrypt-result').innerText;
    await navigator.clipboard.writeText(text);
    alert('Скопировано!');
}

// Закрытие модального окна по клику вне его
document.addEventListener('DOMContentLoaded', () => {
    const modal = document.getElementById('admin-delete-confirmation-modal');
    if (modal) {
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                closeAdminDeleteConfirmation();
            }
        });
    }
});