CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    email VARCHAR (128) UNIQUE NOT NULL,
    password_hash VARCHAR (128) NOT NULL
);
CREATE TABLE IF NOT EXISTS secrets(
    id SERIAL PRIMARY KEY,
    name VARCHAR (255) NOT NULL,
    content BYTEA,
    version UUID DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
    owner_id INTEGER REFERENCES users (id),
    UNIQUE (name, owner_id)
);
