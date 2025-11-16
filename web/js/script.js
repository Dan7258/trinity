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

    const url = `/encode/${algorithm}`;
    const body = JSON.stringify({ message });

    // Получаем токен
    const token = localStorage.getItem('token');

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                ...(token && { 'Authorization': `Bearer ${token}` }) // Добавляем токен, если он есть
            },
            body: body
        });

        await handleResponse(response, output, (data) => {
            let jsonToCopy;

            if (algorithm === 'rsa') {
                jsonToCopy = JSON.stringify({
                    encrypted_message: data.encrypted_message,
                    d: data.d,
                    n: data.n
                }, null, 2);
                output.dataset.rawJson = jsonToCopy;

                output.innerHTML = `
                    <strong>Зашифрованное сообщение:</strong><br>
                    <div style="overflow-x:auto;"><code>${data.encrypted_message || '—'}</code></div><br>
                    <strong>d:</strong><br>
                    <div style="overflow-x:auto;"><code>${data.d || '—'}</code></div><br>
                    <strong>n:</strong><br>
                    <div style="overflow-x:auto;"><code>${data.n || '—'}</code></div>
                `;
            }
            else if (algorithm === 'kuznechik') {
                jsonToCopy = JSON.stringify({
                    encrypted_message: data.encrypted_message,
                    key: data.key
                }, null, 2);
                output.dataset.rawJson = jsonToCopy;

                output.innerHTML = `
                    <strong>Шифротекст:</strong><br>
                    <div style="overflow-x:auto;"><code>${data.encrypted_message || '—'}</code></div><br>
                    <strong>Ключ:</strong><br>
                    <div style="overflow-x:auto;"><code>${data.key || '—'}</code></div>
                `;
            }
            else if (algorithm === 'stribog') {
                const hash = data.encrypted_message || data.hash || '';
                jsonToCopy = JSON.stringify({ hash }, null, 2);
                output.dataset.rawJson = jsonToCopy;
                output.innerText = hash;
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
    }
    else if (algorithm === 'kuznechik') {
        let encryptedObj;
        try {
            encryptedObj = JSON.parse(rawInput);
        } catch (e) {
            output.innerHTML = `<span style="color: red;">Неверный JSON. Ожидается: {"encrypted_message": "...", "key": "..."}</span>`;
            return;
        }

        if (!encryptedObj.encrypted_message || !encryptedObj.key) {
            output.innerHTML = `<span style="color: red;">JSON должен содержать поля: encrypted_message, key</span>`;
            return;
        }

        url = '/decode/kuznechik';
        body = JSON.stringify({
            encrypted_message: encryptedObj.encrypted_message,
            key: encryptedObj.key
        });
    }
    else {
        output.innerHTML = `<span style="color: red;">Расшифровка для ${algorithm} не поддерживается</span>`;
        return;
    }

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: body
        });

        await handleResponse(response, output, (data) => {
            const message = data.message || '';
            const jsonToCopy = JSON.stringify({ message }, null, 2);
            output.dataset.rawJson = jsonToCopy;
            output.innerText = message || 'Расшифровано (пусто)';
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

async function loadHistory() {
    const token = localStorage.getItem('token');
    const historyTab = document.getElementById('history-tab');
    const historyContent = document.getElementById('history-content');

    if (!token) {
        if (historyTab) historyTab.style.display = 'none';
        return;
    } else {
        if (historyTab) historyTab.style.display = 'block';
    }

    if (!historyContent) return;

    historyContent.innerHTML = '<p><em>Загрузка истории...</em></p>';

    const algorithms = ['kuznechik', 'rsa', 'stribog'];
    let allHistory = [];

    try {
        for (const alg of algorithms) {
            const response = await fetch(`/history/${alg}`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            if (!response.ok && response.status !== 404) continue;

            const data = await response.json();

            if (Array.isArray(data)) {
                data.forEach(item => {
                    item.algorithm = alg;
                    item.operation = alg === 'stribog' ? 'hash' : 'encrypt';
                    allHistory.push(item);
                });
            }
        }

        allHistory.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));

        if (allHistory.length === 0) {
            historyContent.innerHTML = '<p style="color: #666;">История пуста.</p>';
            return;
        }

        let html = `
            <div style="overflow-x: auto;">
                <table style="width:100%; border-collapse: collapse; font-size: 13px;">
                    <thead>
                        <tr style="background: #f0f0f0; text-align: left;">
                            <th style="padding: 8px; border-bottom: 2px solid #ddd;">Алгоритм</th>
                            <th style="padding: 8px; border-bottom: 2px solid #ddd;">Время</th>
                            <th style="padding: 8px; border-bottom: 2px solid #ddd;">Результат</th>
                            <th style="padding: 8px; border-bottom: 2px solid #ddd;">Действие</th>
                        </tr>
                    </thead>
                    <tbody>
        `;

        allHistory.forEach(item => {
            // Парсим created_at - он может быть в разных форматах
            let time = 'Invalid Date';
            try {
                if (item.created_at) {
                    const date = new Date(item.created_at);
                    if (!isNaN(date.getTime())) {
                        time = date.toLocaleString('ru-RU');
                    }
                }
            } catch (e) {
                console.error('Date parse error:', e);
            }

            const algName = { kuznechik: 'Кузнечик', rsa: 'RSA', stribog: 'Стрибог' }[item.algorithm];

            const canDecrypt = ['kuznechik', 'rsa'].includes(item.algorithm);
            let decryptPayload = null;
            let preview = '';

            if (item.algorithm === 'kuznechik') {
                // Данные приходят в формате: { encrypted_message: "...", key: "..." }
                const encMsg = item.encrypted_message || '';
                const key = item.key || '';
                preview = encMsg.substring(0, 60) + (encMsg.length > 60 ? '...' : '');
                if (encMsg && key) {
                    decryptPayload = {
                        encrypted_message: encMsg,
                        key: key
                    };
                }
            } else if (item.algorithm === 'rsa') {
                // Данные приходят в формате: { encrypted_message: "...", d: "...", n: "..." }
                const encMsg = item.encrypted_message || '';
                const d = item.d || '';
                const n = item.n || '';
                preview = encMsg.substring(0, 60) + (encMsg.length > 60 ? '...' : '');
                if (encMsg && d && n) {
                    decryptPayload = {
                        encrypted_message: encMsg,
                        d: d,
                        n: n
                    };
                }
            } else if (item.algorithm === 'stribog') {
                // Данные приходят в формате: { hash: "..." }
                const hash = item.hash || '';
                preview = hash.substring(0, 60) + (hash.length > 60 ? '...' : '');
            }

            html += `
                <tr style="border-bottom: 1px solid #eee;">
                    <td style="padding: 8px;">${algName}</td>
                    <td style="padding: 8px; white-space: nowrap;">${time}</td>
                    <td style="padding: 8px; font-family: monospace; font-size: 12px; max-width: 250px; overflow: hidden; text-overflow: ellipsis;" title="${preview}">
                        ${preview || '—'}
                    </td>
                    <td style="padding: 8px;">
                        ${canDecrypt ? `
                            <button class="copy-btn" style="font-size: 11px; padding: 4px 8px; background: #007bff; color: white;" 
                                    onclick='decryptFromHistory("${item.algorithm}", ${JSON.stringify(decryptPayload).replace(/'/g, "&#39;")})'>
                                Расшифровать
                            </button>
                        ` : `<span style="color: #aaa; font-size: 11px;">—</span>`}
                    </td>
                </tr>
            `;
        });

        html += `</tbody></table></div>`;
        historyContent.innerHTML = html;

    } catch (error) {
        historyContent.innerHTML = `<p style="color: red;">Ошибка: ${error.message}</p>`;
    }
}

// === Расшифровка из истории ===
async function decryptFromHistory(algorithm, payload) {
    if (!payload || !algorithm) return;

    const resultEl = document.getElementById('decrypt-result');
    const modal = document.getElementById('decrypt-modal');
    resultEl.innerHTML = '<em>Обработка... Это может занять некоторое время для RSA.</em>';
    modal.style.display = 'flex';

    try {
        // Увеличиваем таймаут для RSA
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), 60000); // 60 секунд

        const response = await fetch(`/decode/${algorithm}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: JSON.stringify(payload),
            signal: controller.signal
        });

        clearTimeout(timeoutId);

        const contentType = response.headers.get("content-type");
        console.log('Content-Type:', contentType); // Отладка
        console.log('Response status:', response.status); // Отладка

        let data;

        // Пытаемся распарсить как JSON, даже если Content-Type неправильный
        try {
            const text = await response.text();
            console.log('Сырой ответ:', text); // Отладка
            data = JSON.parse(text);
            console.log('Распарсенный JSON:', data); // Отладка
        } catch (e) {
            console.error('Ошибка парсинга JSON:', e);
            resultEl.innerHTML = `<span style="color: red;">Ошибка парсинга ответа: ${e.message}</span>`;
            return;
        }

        if (response.status >= 400) {
            resultEl.innerHTML = `<span style="color: red;">Ошибка: ${data.message || 'Неизвестная ошибка'}</span>`;
            return;
        }

        const message = data.message || data.decrypted_message || '(пусто)';
        console.log('Расшифрованное сообщение:', message); // Отладка
        resultEl.innerText = message;

        // Добавляем кнопку копирования
        resultEl.innerHTML += `\n\n<button onclick="copyText('${message.replace(/'/g, "\\'")}'); event.stopPropagation();" style="margin-top:10px; padding:6px 12px; font-size:12px;">Скопировать результат</button>`;

    } catch (error) {
        if (error.name === 'AbortError') {
            resultEl.innerHTML = `<span style="color: red;">Таймаут: операция заняла слишком много времени (более 60 секунд)</span>`;
        } else {
            resultEl.innerHTML = `<span style="color: red;">Сетевая ошибка: ${error.message}</span>`;
        }
        console.error('Ошибка расшифровки:', error);
    }
}

// === Копирование произвольного текста (для истории) ===
async function copyText(text) {
    try {
        await navigator.clipboard.writeText(text);
        alert('Скопировано в буфер обмена!');
    } catch (err) {
        alert('Ошибка копирования: ' + err.message);
    }
}

function closeDecryptModal() {
    document.getElementById('decrypt-modal').style.display = 'none';
}

// === Перезагрузка истории при открытии вкладки ===
document.querySelector('[data-tab="history"]')?.addEventListener('click', () => {
    setTimeout(loadHistory, 100); // небольшая задержка для активации вкладки
});

// === Автозагрузка при старте, если вкладка активна или token есть ===
document.addEventListener('DOMContentLoaded', () => {
    // Проверяем token и показываем вкладку
    const token = localStorage.getItem('token');
    const historyTab = document.getElementById('history-tab');
    if (historyTab) {
        historyTab.style.display = token ? 'block' : 'none';
    }

    // Если открыта вкладка истории — загружаем
    if (document.querySelector('.tab-button.active')?.dataset.tab === 'history') {
        loadHistory();
    }
});

// === Обновление при изменении token (через storage) ===
window.addEventListener('storage', (e) => {
    if (e.key === 'token') {
        const historyTab = document.getElementById('history-tab');
        if (historyTab) {
            historyTab.style.display = e.newValue ? 'block' : 'none';
        }
        if (e.newValue && document.getElementById('history').classList.contains('active')) {
            loadHistory();
        }
    }
});