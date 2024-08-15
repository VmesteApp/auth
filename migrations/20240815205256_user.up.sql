CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    email VARCHAR(255),
    passHash VARCHAR(255)
);