CREATE TABLE
    IF NOT EXISTS users (
        id serial PRIMARY KEY,
        email VARCHAR(255) NULL UNIQUE,
        pass_hash VARCHAR(255) NULL
    );