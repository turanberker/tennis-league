CREATE TABLE tennisleague.matche (
    id varchar(100) PRIMARY KEY DEFAULT gen_random_uuid(),

    league_id varchar(100) ,
    team_1_id varchar(100) NOT NULL,
    team_2_id varchar(100) NOT NULL,

    team_1_score INT,
    team_2_score INT,
    winner_id varchar(100),

    match_date TIMESTAMPTZ NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    approve_date  TIMESTAMPTZ NULL,
    CONSTRAINT fk_matches_league
        FOREIGN KEY (league_id)
        REFERENCES tennisleague.league(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_matches_team1
        FOREIGN KEY (team_1_id)
        REFERENCES tennisleague.team(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_matches_team2
        FOREIGN KEY (team_2_id)
        REFERENCES tennisleague.team(id)
        ON DELETE CASCADE,

     CONSTRAINT fk_matches_team_winner
            FOREIGN KEY (winner_id)
            REFERENCES tennisleague.team(id)
            ON DELETE CASCADE,

    CONSTRAINT chk_different_teams
        CHECK (team_1_id <> team_2_id)
);

CREATE TABLE tennisleague.match_set (
    id varchar(100) PRIMARY KEY DEFAULT gen_random_uuid(),

    match_id varchar(100) NOT NULL,
    set_number INT NOT NULL,          -- 1,2,3...

    team_1_games INT ,
    team_2_games INT ,

    team_1_tie_break_score INT ,
    team_2_tie_break_score INT,

    CONSTRAINT fk_set_match
        FOREIGN KEY (match_id)
        REFERENCES tennisleague.matche(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_fixture_set UNIQUE (match_id, set_number),

    CONSTRAINT chk_games_positive
        CHECK (team_1_games >= 0 AND team_2_games >= 0)
);