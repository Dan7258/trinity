// Переключение вкладок
document.querySelectorAll('.tab-button').forEach(button => {
    button.addEventListener('click', () => {
        document.querySelectorAll('.tab-button').forEach(btn => btn.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));

        button.classList.add('active');
        document.getElementById(button.dataset.tab).classList.add('active');
    });
});

// Шифрование / Хэширование
async function encrypt(algorithm) {
    const input = document.getElementById(`${algorithm}-input`);
    const output = document.getElementById(`${algorithm}-encrypt-output`);
    const message = input.value.trim();

    if (!message) {
        alert('Введите сообщение');
        return;
    }

    output.innerText = 'Обработка...';

    try {
        const response = await fetch(`/encrypt/${algorithm}`, {  // ← Замени путь
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ message })
        });

        const data = await response.json();
        output.innerText = data.cipher || data.hash || JSON.stringify(data, null, 2);
    } catch (error) {
        output.innerText = 'Ошибка: ' + error.message;
    }
}

// Расшифровка (только для kuznechik и rsa)
async function decrypt(algorithm) {
    const input = document.getElementById(`${algorithm}-decrypt-input`);
    const output = document.getElementById(`${algorithm}-decrypt-output`);
    const cipher = input.value.trim();

    if (!cipher) {
        alert('Введите зашифрованное сообщение');
        return;
    }

    output.innerText = 'Обработка...';

    try {
        const response = await fetch(`/decrypt/${algorithm}`, {  // ← Замени путь
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ cipher })
        });

        const data = await response.json();
        output.innerText = data.message || JSON.stringify(data, null, 2);
    } catch (error) {
        output.innerText = 'Ошибка: ' + error.message;
    }
}