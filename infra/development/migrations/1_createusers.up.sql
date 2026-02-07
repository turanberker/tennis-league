CREATE TABLE tennisleague.users (
    id BIGSERIAL PRIMARY KEY,
    email  VARCHAR(100) NOT NULL UNIQUE,
    phone  VARCHAR(100),
    name  VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    password_hash TEXT NOT NULL,
    role  VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    approved BOOLEAN DEFAULT FALSE    
);