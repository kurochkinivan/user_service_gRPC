-- Active: 1744705463730@@127.0.0.1@5433@users
CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid() NOT NULL,
    name TEXT NOT NULL,
    age INT NOT NULL,
    gender TEXT NOT NULL,
    about TEXT NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT check_gender CHECK (gender IN ('male', 'female'))
);