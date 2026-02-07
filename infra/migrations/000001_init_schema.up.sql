

CREATE TABLE tennisleague.leagues (
    id BIGSERIAL PRIMARY KEY,
    name  VARCHAR(100) NOT NULL
);

CREATE TABLE tennisleague.teams (
    id BIGSERIAL PRIMARY KEY,
    league_id BIGINT NOT NULL REFERENCES leagues(id) ON DELETE CASCADE,
    name  VARCHAR(100) NOT NULL
);

CREATE TABLE tennisleague.team_players (
    team_id BIGINT REFERENCES teams(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (team_id, user_id)
);

/* Bir oyuncu sadece 1 takımda olabilir */
CREATE UNIQUE INDEX ux_team_players_user
ON team_players(user_id);

