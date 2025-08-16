-- вариант через staging с меньшей нагрузкой на WAL
CREATE UNLOGGED TABLE users_stage (LIKE users);
COPY users_stage (id,email,password_hash,salt,first_name,last_name,birth_date,created_at,is_active,expires_at)
FROM '/docker-entrypoint-initdb.d/users.csv'
WITH (FORMAT csv, HEADER false, DELIMITER ',', QUOTE '"');

-- валидации/чистки при необходимости
INSERT INTO users
SELECT * FROM users_stage
ON CONFLICT (id) DO NOTHING;

DROP TABLE users_stage;
ANALYZE users;
