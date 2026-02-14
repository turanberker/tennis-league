CREATE TABLE tennisleague.players (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    user_id BIGINT UNIQUE,
    CONSTRAINT fk_player_user
        FOREIGN KEY (user_id)
        REFERENCES tennisleague.users(id)
        ON DELETE SET NULL
);

CREATE UNIQUE INDEX uq_player_uuid
ON tennisleague.players (uuid);