// static/js/login.js
// Работает с POST /secure/two-fa/{email} → 200 OK + пустое тело

document.addEventListener('DOMContentLoaded', () => {
    const form          = document.getElementById('login-form');
    const emailInput    = document.getElementById('login-email');
    const passwordInput = document.getElementById('login-password');
    const codeGroup     = document.getElementById('code-group');
    const codeInput     = document.getElementById('login-code');
    const submitBtn     = document.getElementById('submit-btn');
    const msg           = document.getElementById('login-message');

    let twoFaSent = false;

    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const email    = emailInput.value.trim();
        const password = passwordInput.value;

        // === Первый клик — запрос 2FA ===
        if (!twoFaSent) {
            if (!email || !password) return;

            // простая проверка формата email
            if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
                msg.style.color = 'red';
                msg.textContent = 'Введите корректный email';
                return;
            }

            submitBtn.disabled = true;
            submitBtn.textContent = 'Отправка кода...';
            msg.textContent = '';
            msg.style.color = '#007bff';

            try {
                const res = await fetch(`/secure/two-fa/${encodeURIComponent(email)}`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' }
                });

                // Главное — статус 200 → считаем успехом (даже если тело пустое)
                if (res.ok) {
                    twoFaSent = true;
                    codeGroup.classList.remove('hidden');
                    codeInput.focus();

                    msg.style.color = 'green';
                    msg.textContent = 'Код отправлен на вашу почту';

                    submitBtn.textContent = 'Войти';
                } else {
                    // 4xx / 5xx
                    const text = await res.text();
                    msg.style.color = 'red';
                    msg.textContent = text || 'Не удалось отправить код';
                    submitBtn.textContent = 'Войти';
                }
            } catch (err) {
                console.error(err);
                msg.style.color = 'red';
                msg.textContent = 'Сетевая ошибка';
                submitBtn.textContent = 'Войти';
            } finally {
                submitBtn.disabled = false;
            }

            return;
        }

        // === Второй клик — финальный логин ===
        const code = codeInput.value.trim();
        if (!code) {
            msg.style.color = 'red';
            msg.textContent = 'Введите код подтверждения';
            return;
        }

        submitBtn.disabled = true;
        submitBtn.textContent = 'Вход...';
        msg.textContent = '';

        try {
            const res = await fetch('/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password, code })
            });

            const data = await res.json();

            if (res.ok && data.token) {
                localStorage.setItem('token', data.token);
                msg.style.color = 'green';
                msg.textContent = 'Успешный вход!';

                // обновляем хедер во всех вкладках
                window.dispatchEvent(new StorageEvent('storage', { key: 'token' }));

                setTimeout(() => window.location.href = '/index.html', 800);
            } else {
                msg.style.color = 'red';
                msg.textContent = data.message || 'Неверный логин, пароль или код';
            }
        } catch (err) {
            console.error(err);
            msg.style.color = 'red';
            msg.textContent = 'Сетевая ошибка';
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = 'Войти';
        }
    });
});