// static/js/delete-account.js

// Показать/скрыть секцию удаления аккаунта в зависимости от авторизации
function updateDeleteAccountVisibility() {
    const token = localStorage.getItem('token');
    const deleteSection = document.getElementById('delete-account-section');

    if (deleteSection) {
        deleteSection.style.display = token ? 'block' : 'none';
    }
}

// Показать модальное окно подтверждения
function showDeleteConfirmation() {
    const modal = document.getElementById('delete-confirmation-modal');
    if (modal) {
        modal.style.display = 'flex';
    }
}

// Закрыть модальное окно
function closeDeleteConfirmation() {
    const modal = document.getElementById('delete-confirmation-modal');
    if (modal) {
        modal.style.display = 'none';
    }
}


// Удалить аккаунт
async function deleteAccount() {
    const token = localStorage.getItem('token');
    if (!token) {
        alert('Требуется авторизация');
        window.location.href = '/login.html';
        return;
    }

    try {
        const response = await fetch(`/api/delete-user`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        });

        if (response.ok) {
            alert('Аккаунт успешно удалён');

            // Очищаем данные
            localStorage.removeItem('token');
            document.cookie = 'token=; path=/; max-age=0';

            // Закрываем модалку и перенаправляем
            closeDeleteConfirmation();
            window.location.href = '/register.html';

        } else if (response.status === 401) {
            alert('Сессия истекла. Войдите заново.');
            localStorage.removeItem('token');
            window.location.href = '/login.html';

        } else if (response.status === 403) {
            alert('Доступ запрещён. Вы можете удалить только свой аккаунт.');
            closeDeleteConfirmation();

        } else if (response.status === 404) {
            alert('Пользователь не найден');
            closeDeleteConfirmation();

        } else {
            const data = await response.json().catch(() => ({}));
            alert(`Ошибка: ${data.message || 'Не удалось удалить аккаунт'}`);
            closeDeleteConfirmation();
        }

    } catch (err) {
        console.error('Ошибка удаления аккаунта:', err);
        alert('Сетевая ошибка: ' + err.message);
        closeDeleteConfirmation();
    }
}

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    updateDeleteAccountVisibility();
});

// Обновление при изменении токена
window.addEventListener('storage', (e) => {
    if (e.key === 'token') {
        updateDeleteAccountVisibility();
    }
});

// Закрытие модалки по клику вне её
document.addEventListener('DOMContentLoaded', () => {
    const modal = document.getElementById('delete-confirmation-modal');
    if (modal) {
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                closeDeleteConfirmation();
            }
        });
    }
});