CREATE TABLE tennisleague.player (
    id VARCHAR(100) PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    sex VARCHAR(1) NOT NULL,
    user_id VARCHAR(100) UNIQUE,
    CONSTRAINT fk_player_user
        FOREIGN KEY (user_id)
        REFERENCES tennisleague.user(id)
        ON DELETE SET NULL
);

CREATE  INDEX sex
ON tennisleague.player (sex);