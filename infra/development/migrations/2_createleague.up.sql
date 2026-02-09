CREATE TABLE tennisleague.leagues (
    id BIGSERIAL PRIMARY KEY,
    name  VARCHAR(100) NOT NULL unique
);