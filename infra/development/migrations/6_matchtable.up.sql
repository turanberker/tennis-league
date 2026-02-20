CREATE TABLE tennisleague.matches (
    id varchar(100) PRIMARY KEY DEFAULT gen_random_uuid(),

    league_id varchar(100) ,
    team_1_id varchar(100) NOT NULL,
    team_2_id varchar(100) NOT NULL,

    match_date TIMESTAMPTZ NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',

       CONSTRAINT fk_matches_league
        FOREIGN KEY (league_id)
        REFERENCES tennisleague.leagues(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_matches_team1
        FOREIGN KEY (team_1_id)
        REFERENCES tennisleague.teams(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_matches_team2
        FOREIGN KEY (team_2_id)
        REFERENCES tennisleague.teams(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_different_teams
        CHECK (team_1_id <> team_2_id)
);

CREATE TABLE tennisleague.match_sets (
    id varchar(100) PRIMARY KEY DEFAULT gen_random_uuid(),

    match_id varchar(100) NOT NULL,
    set_number INT NOT NULL,          -- 1,2,3...

     team_1_games INT ,
    team_2_games INT ,

   team_1_tie_break_score INT ,
    team_2_tie_break_score INT,

    CONSTRAINT fk_set_match
        FOREIGN KEY (match_id)
        REFERENCES tennisleague.matches(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_fixture_set UNIQUE (match_id, set_number),

    CONSTRAINT chk_games_positive
        CHECK (team_1_games >= 0 AND team_2_games >= 0)
);