CREATE TABLE tennisleague.league_coordinator (
    league_id  VARCHAR(100) Not NULL,
    user_id  VARCHAR(100) NOT NULL,

    CONSTRAINT fk_league_coordinator_league FOREIGN KEY (league_id) REFERENCES tennisleague.league(id),
    CONSTRAINT fk_league_coordinator_user FOREIGN KEY (user_id) REFERENCES tennisleague.user(id),

    -- İkisinin kombinasyonunu benzersiz ve birincil anahtar yapar
    PRIMARY KEY (league_id, user_id)
);

CREATE INDEX idx_league_coordinator_league ON tennisleague.league_coordinator (league_id);
CREATE INDEX idx_league_coordinator_user ON tennisleague.league_coordinator (user_id);