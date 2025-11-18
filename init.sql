
-- Пользователи
CREATE TABLE IF NOT EXISTS "users" (
                                      id            SERIAL PRIMARY KEY,
                                      login         VARCHAR(255) NOT NULL UNIQUE,
                                      role          VARCHAR(255) NOT NULL,
                                      password      VARCHAR(255) NOT NULL
);

-- --------------------------------------------------------------
--  Кузнечик (симметричное шифрование)
--  EncryptedData { EncryptedMessage, Key }
-- --------------------------------------------------------------
CREATE TABLE IF NOT EXISTS kuznechiks (
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
CREATE TABLE IF NOT EXISTS stribogs (
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
CREATE TABLE IF NOT EXISTS rsas (
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
COMMENT ON TABLE kuznechiks IS 'Данные, зашифрованные алгоритмом Кузнечик (с ключом)';
COMMENT ON TABLE stribogs   IS 'Хэши, вычисленные по алгоритму Стрибог';
COMMENT ON TABLE rsas       IS 'Данные, зашифрованные RSA (с закрытым ключом d, n)';
