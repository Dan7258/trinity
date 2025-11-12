// === Переключение вкладок ===
document.querySelectorAll('.tab-button').forEach(button => {
    button.addEventListener('click', () => {
        document.querySelectorAll('.tab-button').forEach(btn => btn.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));

        button.classList.add('active');
        document.getElementById(button.dataset.tab).classList.add('active');
    });
});

// === Универсальная обработка ответа сервера ===
async function handleResponse(response, outputElement, successCallback) {
    const contentType = response.headers.get("content-type");
    if (!contentType || !contentType.includes("application/json")) {
        const text = await response.text();
        outputElement.innerHTML = `<span style="color: red;">Ошибка сервера: ${response.status} ${text}</span>`;
        return;
    }

    const data = await response.json();

    if (response.status >= 400) {
        outputElement.innerHTML = `<span style="color: red;">Ошибка: ${data.message || 'Неизвестная ошибка'}</span>`;
        return;
    }

    if (response.status >= 200 && response.status < 300) {
        successCallback(data);
    } else {
        outputElement.innerHTML = `<span style="color: orange;">Неожиданный статус: ${response.status}</span>`;
    }
}

// === Шифрование / Хэширование ===
async function encrypt(algorithm) {
    const input = document.getElementById(`${algorithm}-input`);
    const output = document.getElementById(`${algorithm}-encrypt-output`);
    const message = input.value.trim();

    if (!message) {
        alert('Введите сообщение');
        return;
    }

    output.innerHTML = '<em>Обработка...</em>';

    let url, body;

    if (algorithm === 'rsa') {
        url = '/encode/rsa';
        body = JSON.stringify({ message });
    } else if (algorithm === 'stribog') {
        url = '/encrypt/stribog';  // ← Замени, если путь другой
        body = JSON.stringify({ message });
    } else {
        url = `/encrypt/${algorithm}`;
        body = JSON.stringify({ message });
    }

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: body
        });

        await handleResponse(response, output, (data) => {
            let jsonToCopy;

            if (algorithm === 'rsa') {
                jsonToCopy = JSON.stringify(data, null, 2);
                output.dataset.rawJson = jsonToCopy;

                output.innerHTML = `
                    <strong>Зашифрованное сообщение:</strong><br>
                    <div style="overflow-x:auto; max-width:100%;"><code>${data.encrypted_message || '—'}</code></div><br>
                    <strong>d:</strong><br>
                    <div style="overflow-x:auto; max-width:100%;"><code>${data.d || '—'}</code></div><br>
                    <strong>n:</strong><br>
                    <div style="overflow-x:auto; max-width:100%;"><code>${data.n || '—'}</code></div><br><br>
                    <em>Нажмите "Копировать результат" ниже:</em>
                `;
            } else if (algorithm === 'stribog') {
                const result = data.hash || JSON.stringify(data, null, 2);
                jsonToCopy = JSON.stringify({ hash: data.hash }, null, 2);
                output.dataset.rawJson = jsonToCopy;
                output.innerText = result;
            } else {
                const result = data.cipher || JSON.stringify(data, null, 2);
                jsonToCopy = JSON.stringify({ cipher: data.cipher }, null, 2);
                output.dataset.rawJson = jsonToCopy;
                output.innerText = result;
            }
        });

    } catch (error) {
        output.innerHTML = `<span style="color: red;">Сетевая ошибка: ${error.message}</span>`;
    }
}

// === Расшифровка (только kuznechik и rsa) ===
async function decrypt(algorithm) {
    const input = document.getElementById(`${algorithm}-decrypt-input`);
    const output = document.getElementById(`${algorithm}-decrypt-output`);
    const rawInput = input.value.trim();

    if (!rawInput) {
        alert('Введите данные для расшифровки');
        return;
    }

    output.innerHTML = '<em>Обработка...</em>';

    let url, body;

    if (algorithm === 'rsa') {
        let encryptedObj;
        try {
            encryptedObj = JSON.parse(rawInput);
        } catch (e) {
            output.innerHTML = `<span style="color: red;">Неверный JSON. Ожидается: {"encrypted_message": "...", "d": "...", "n": "..."}</span>`;
            return;
        }

        if (!encryptedObj.encrypted_message || !encryptedObj.d || !encryptedObj.n) {
            output.innerHTML = `<span style="color: red;">JSON должен содержать поля: encrypted_message, d, n</span>`;
            return;
        }

        url = '/decode/rsa';
        body = JSON.stringify({
            encrypted_message: encryptedObj.encrypted_message,
            d: encryptedObj.d,
            n: encryptedObj.n
        });
    } else {
        url = `/decrypt/${algorithm}`;
        body = JSON.stringify({ cipher: rawInput });
    }

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: body
        });

        await handleResponse(response, output, (data) => {
            let jsonToCopy;

            if (algorithm === 'rsa') {
                const message = data.message || '';
                jsonToCopy = JSON.stringify({ message }, null, 2);
                output.dataset.rawJson = jsonToCopy;
                output.innerText = message || 'Расшифровано (пусто)';
            } else {
                const message = data.message || JSON.stringify(data, null, 2);
                jsonToCopy = JSON.stringify({ message: data.message }, null, 2);
                output.dataset.rawJson = jsonToCopy;
                output.innerText = message;
            }
        });

    } catch (error) {
        output.innerHTML = `<span style="color: red;">Сетевая ошибка: ${error.message}</span>`;
    }
}

// === Копирование результата в буфер обмена ===
async function copyOutput(outputId) {
    const output = document.getElementById(outputId);
    if (!output) return;

    let textToCopy = output.dataset.rawJson;

    if (!textToCopy) {
        textToCopy = output.innerText || output.textContent || '';
        textToCopy = textToCopy.trim();
        if (!textToCopy) {
            alert('Нет данных для копирования');
            return;
        }
        textToCopy = JSON.stringify({ result: textToCopy }, null, 2);
    }

    try {
        await navigator.clipboard.writeText(textToCopy);

        const btn = output.nextElementSibling;
        if (btn && btn.classList.contains('copy-btn')) {
            const originalHTML = btn.innerHTML;
            btn.innerHTML = `
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" style="width:14px;height:14px;">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                Готово
            `;
            btn.style.background = '#1e7e34';

            setTimeout(() => {
                btn.innerHTML = originalHTML;
                btn.style.background = '#28a745';
            }, 1500);
        }
    } catch (err) {
        alert('Ошибка копирования: ' + err.message);
    }
}