// static/js/login.js
// Двухфакторная аутентификация: код с email + код из Telegram бота

document.addEventListener('DOMContentLoaded', () => {
    const form              = document.getElementById('login-form');
    const emailInput        = document.getElementById('login-email');
    const passwordInput     = document.getElementById('login-password');
    const codesGroup        = document.getElementById('codes-group');
    const codeEmailInput    = document.getElementById('login-code-email');
    const codeTelegramInput = document.getElementById('login-code-telegram');
    const submitBtn         = document.getElementById('submit-btn');
    const msg               = document.getElementById('login-message');

    let twoFaSent = false;

    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const email    = emailInput.value.trim();
        const password = passwordInput.value;

        // === Первый клик — запрос 2FA (отправка кода на email) ===
        if (!twoFaSent) {
            if (!email || !password) return;

            // Простая проверка формата email
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

                if (res.ok) {
                    twoFaSent = true;
                    codesGroup.classList.remove('hidden');
                    codeEmailInput.focus();

                    msg.style.color = 'green';
                    msg.innerHTML = 'Код отправлен на почту.<br>Откройте бота <a href="https://t.me/uga_bugaaa_bot" target="_blank" style="color:#0088cc;">@uga_bugaaa_bot</a> для получения второго кода';

                    submitBtn.textContent = 'Войти';
                } else {
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

        // === Второй клик — финальный логин с двумя кодами ===
        const codeEmail    = codeEmailInput.value.trim();
        const codeTelegram = codeTelegramInput.value.trim();

        if (!codeEmail) {
            msg.style.color = 'red';
            msg.textContent = 'Введите код из письма';
            codeEmailInput.focus();
            return;
        }

        if (!codeTelegram) {
            msg.style.color = 'red';
            msg.textContent = 'Введите код из Telegram бота';
            codeTelegramInput.focus();
            return;
        }

        submitBtn.disabled = true;
        submitBtn.textContent = 'Вход...';
        msg.textContent = '';

        try {
            const res = await fetch('/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    email,
                    password,
                    code_email: codeEmail,
                    code_telegram: codeTelegram
                })
            });

            const data = await res.json();

            if (res.ok && data.token) {
                localStorage.setItem('token', data.token);
                msg.style.color = 'green';
                msg.textContent = 'Успешный вход!';

                // Обновляем хедер во всех вкладках
                window.dispatchEvent(new StorageEvent('storage', { key: 'token' }));

                setTimeout(() => window.location.href = '/index.html', 800);
            } else {
                msg.style.color = 'red';
                msg.textContent = data.message || 'Неверный логин, пароль или коды подтверждения';
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