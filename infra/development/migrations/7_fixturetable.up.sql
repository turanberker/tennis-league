CREATE TABLE tennisleague.score_board (
    id varchar(100) PRIMARY KEY DEFAULT gen_random_uuid(),
    league_id varchar(100) not null,
    team_id  varchar(100) not null,

    played INT default 0,
    won INT default 0,
    lost INT default 0,

    won_sets INT default 0,
    lost_sets INT default 0,

    won_games INT default 0,
    lost_games INT default 0,

    score INT default 0,

    CONSTRAINT unique_league_team UNIQUE (league_id, team_id),

    CONSTRAINT fk_scorboard_league
            FOREIGN KEY (league_id)
            REFERENCES tennisleague.leagues(id)
            ON DELETE CASCADE,

    CONSTRAINT fk_scorboard_team
            FOREIGN KEY (team_id)
            REFERENCES tennisleague.teams(id)
            ON DELETE CASCADE
);

CREATE INDEX idx_scorboard_league ON tennisleague.score_board(league_id)