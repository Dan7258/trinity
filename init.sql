
-- Пользователи
CREATE TABLE IF NOT EXISTS "users" (
                                      id            SERIAL PRIMARY KEY,
                                      login         VARCHAR(255) NOT NULL UNIQUE,
                                      password      VARCHAR(255) NOT NULL
);

-- --------------------------------------------------------------
--  Кузнечик (симметричное шифрование)
--  EncryptedData { EncryptedMessage, Key }
-- --------------------------------------------------------------
CREATE TABLE IF NOT EXISTS kuznechik (
                                         id                SERIAL PRIMARY KEY,
                                         user_id           INTEGER NOT NULL
                                             REFERENCES "users"(id) ON DELETE CASCADE,
                                         encrypted_message TEXT NOT NULL,
                                         key               TEXT NOT NULL,
                                         created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- --------------------------------------------------------------
--  Стрибог (хеш)
--  EncryptedData { EncryptedMessage }
-- --------------------------------------------------------------
CREATE TABLE IF NOT EXISTS stribog (
                                       id                SERIAL PRIMARY KEY,
                                       user_id           INTEGER NOT NULL
                                           REFERENCES "users"(id) ON DELETE CASCADE,
                                       encrypted_message TEXT NOT NULL,
                                       created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- --------------------------------------------------------------
--  RSA (асимметричное шифрование)
--  EncryptedData { EncryptedMessage, D, N }
-- --------------------------------------------------------------
CREATE TABLE IF NOT EXISTS rsa (
                                   id                SERIAL PRIMARY KEY,
                                   user_id           INTEGER NOT NULL
                                       REFERENCES "users"(id) ON DELETE CASCADE,
                                   encrypted_message TEXT NOT NULL,
                                   d                 TEXT NOT NULL,
                                   n                 TEXT NOT NULL,
                                   created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- ==============================================================
--  Комментарии
-- ==============================================================
COMMENT ON TABLE "users"   IS 'Пользователи системы';
COMMENT ON TABLE kuznechik IS 'Данные, зашифрованные алгоритмом Кузнечик (с ключом)';
COMMENT ON TABLE stribog   IS 'Хэши, вычисленные по алгоритму Стрибог';
COMMENT ON TABLE rsa       IS 'Данные, зашифрованные RSA (с закрытым ключом d, n)';
