-- Создание таблицы SocialLogin
CREATE TABLE
  IF NOT EXISTS social_logins (
    id serial PRIMARY KEY,
    user_id INT NOT NULL,
    provider VARCHAR(255) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
  );