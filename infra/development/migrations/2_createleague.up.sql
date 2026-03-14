CREATE TABLE tennisleague.league (
    id  VARCHAR(100) PRIMARY KEY  DEFAULT gen_random_uuid(),
    name  VARCHAR(100) NOT NULL unique,
    fixture_created_date TIMESTAMP
);