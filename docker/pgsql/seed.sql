CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    first_name VARCHAR(255) NULL,
    second_name VARCHAR(255) NULL,
    birthdate DATE NULL,
    biography TEXT NULL,
    city VARCHAR(255) NULL
);