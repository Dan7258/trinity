// === Переключение вкладок ===
document.querySelectorAll('.tab-button').forEach(button => {
    button.addEventListener('click', () => {
        document.querySelectorAll('.tab-button').forEach(btn => btn.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));

        button.classList.add('active');
        const tab = button.dataset.tab;
        document.getElementById(tab).classList.add('active');

        // Загружаем историю при открытии вкладки
        loadHistory(tab);
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
    const token = localStorage.getItem('token');

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                ...(token && { 'Authorization': `Bearer ${token}` })
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

            // Перезагружаем историю после успешной операции
            if (token) {
                setTimeout(() => loadHistory(algorithm), 500);
            }
        });

    } catch (error) {
        output.innerHTML = `<span style="color: red;">Сетевая ошибка: ${error.message}</span>`;
    }
}

// === Расшифровка ===
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

// === Копирование результата ===
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

// === Загрузка истории для конкретного алгоритма ===
async function loadHistory(algorithm) {
    const token = localStorage.getItem('token');
    const historySection = document.getElementById(`${algorithm}-history-section`);
    const historyContent = document.getElementById(`${algorithm}-history-content`);

    if (!token || !historySection || !historyContent) {
        if (historySection) historySection.style.display = 'none';
        return;
    }

    historySection.style.display = 'block';
    historyContent.innerHTML = '<p><em>Загрузка истории...</em></p>';

    try {
        const response = await fetch(`/api/history/${algorithm}`, {
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        });

        if (!response.ok) {
            if (response.status === 404) {
                historyContent.innerHTML = '<p style="color: #666;">История пуста.</p>';
                return;
            }
            throw new Error(`HTTP ${response.status}`);
        }

        const data = await response.json();

        if (!Array.isArray(data) || data.length === 0) {
            historyContent.innerHTML = '<p style="color: #666;">История пуста.</p>';
            return;
        }

        // Сортируем по времени (новые первые)
        data.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));

        let html = `
            <table class="history-table">
                <thead>
                    <tr>
                        <th>Время</th>
                        <th style="width: 100%; min-width: 400px;">Результат</th>
                        <th>Действие</th>
                    </tr>
                </thead>
                <tbody>
        `;

        data.forEach(item => {
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

            const canDecrypt = ['kuznechik', 'rsa'].includes(algorithm);
            let decryptPayload = null;
            let fullText = '';
            let preview = '';

            if (algorithm === 'kuznechik') {
                const encMsg = item.encrypted_message || '';
                const key = item.key || '';
                fullText = `Шифротекст: ${encMsg}\nКлюч: ${key}`;
                preview = encMsg.substring(0, 60) + (encMsg.length > 60 ? '...' : '');
                if (encMsg && key) {
                    decryptPayload = { encrypted_message: encMsg, key: key };
                }
            } else if (algorithm === 'rsa') {
                const encMsg = item.encrypted_message || '';
                const d = item.d || '';
                const n = item.n || '';
                fullText = `Зашифровано: ${encMsg}\nd: ${d}\nn: ${n}`;
                preview = encMsg.substring(0, 60) + (encMsg.length > 60 ? '...' : '');
                if (encMsg && d && n) {
                    decryptPayload = { encrypted_message: encMsg, d: d, n: n };
                }
            } else if (algorithm === 'stribog') {
                fullText = item.encrypted_message || '';
                preview = fullText.substring(0, 60) + (fullText.length > 60 ? '...' : '');
            }

            // Кликабельная ячейка — раскрывает полный текст
            html += `
                <tr>
                    <td style="white-space: nowrap;">${time}</td>
                    <td class="full-text-cell" title="Кликните, чтобы скопировать полный результат">
                        <div class="text-preview">${preview || '—'}</div>
                        <div class="text-full" style="display: none; margin-top: 8px; padding: 8px; background: #f8f9fa; border-radius: 4px; white-space: pre-wrap; word-break: break-all; font-family: monospace; font-size: 12px;">
                            ${fullText || '—'}
                        </div>
                    </td>
                    <td>
                        ${canDecrypt && decryptPayload ? `
                            <button class="copy-btn" style="font-size: 11px; padding: 4px 8px; background: #007bff; color: white;" 
                                    onclick='decryptFromHistory("${algorithm}", ${JSON.stringify(decryptPayload).replace(/'/g, "&#39;")})'>
                                Расшифровать
                            </button>
                        ` : `<span style="color: #aaa; font-size: 11px;">—</span>`}
                    </td>
                </tr>
            `;
        });

        html += `</tbody></table>`;
        historyContent.innerHTML = html;

        // Добавляем обработчик клика по ячейкам с полным текстом
        document.querySelectorAll('.full-text-cell').forEach(cell => {
            cell.addEventListener('click', function() {
                const full = this.querySelector('.text-full');
                const preview = this.querySelector('.text-preview');
                if (full.style.display === 'none') {
                    full.style.display = 'block';
                    preview.style.display = 'none';
                } else {
                    full.style.display = 'none';
                    preview.style.display = 'block';
                }
            });

            // Двойной клик — копирует полный текст
            cell.addEventListener('dblclick', function(e) {
                e.stopPropagation();
                const text = this.querySelector('.text-full').innerText;
                navigator.clipboard.writeText(text).then(() => {
                    alert('Полный результат скопирован в буфер обмена!');
                }).catch(() => {
                    alert('Не удалось скопировать');
                });
            });
        });

    } catch (error) {
        historyContent.innerHTML = `<p style="color: red;">Ошибка загрузки: ${error.message}</p>`;
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
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), 60000);

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

        let data;
        try {
            const text = await response.text();
            data = JSON.parse(text);
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
        resultEl.innerText = message;

    } catch (error) {
        if (error.name === 'AbortError') {
            resultEl.innerHTML = `<span style="color: red;">Таймаут: операция заняла слишком много времени (более 60 секунд)</span>`;
        } else {
            resultEl.innerHTML = `<span style="color: red;">Сетевая ошибка: ${error.message}</span>`;
        }
        console.error('Ошибка расшифровки:', error);
    }
}

// === Закрытие модального окна ===
function closeDecryptModal() {
    document.getElementById('decrypt-modal').style.display = 'none';
}

// === Копирование результата из модального окна ===
async function copyModalResult() {
    const resultEl = document.getElementById('decrypt-result');
    const text = resultEl.innerText || resultEl.textContent || '';

    if (!text.trim()) {
        alert('Нет данных для копирования');
        return;
    }

    try {
        await navigator.clipboard.writeText(text);
        alert('Скопировано в буфер обмена!');
    } catch (err) {
        alert('Ошибка копирования: ' + err.message);
    }
}

// === Инициализация при загрузке страницы ===
document.addEventListener('DOMContentLoaded', () => {
    const token = localStorage.getItem('token');

    // Загружаем историю для активной вкладки, если пользователь авторизован
    if (token) {
        const activeTab = document.querySelector('.tab-button.active');
        if (activeTab) {
            loadHistory(activeTab.dataset.tab);
        }
    }
});

// === Обновление при изменении token в localStorage ===
window.addEventListener('storage', (e) => {
    if (e.key === 'token') {
        if (e.newValue) {
            const activeTab = document.querySelector('.tab-button.active');
            if (activeTab) {
                loadHistory(activeTab.dataset.tab);
            }
        } else {
            document.querySelectorAll('.history-section').forEach(section => {
                section.style.display = 'none';
            });
        }
    }
});