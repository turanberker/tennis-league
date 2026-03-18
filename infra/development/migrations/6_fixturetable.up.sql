CREATE TABLE tennisleague.score_board (
    id varchar(100) PRIMARY KEY DEFAULT gen_random_uuid(),
    league_id varchar(100) not null,

    -- Katılımcı: Ya oyuncu ya takım dolu olmalı
    player_id VARCHAR(100) NULL,
    team_id   VARCHAR(100) NULL,

    played INT default 0,
    won INT default 0,
    lost INT default 0,

    won_sets INT default 0,
    lost_sets INT default 0,

    won_games INT default 0,
    lost_games INT default 0,

    score INT default 0,

    -- Sanal Sütun: Set Averajı
    set_diff INTEGER GENERATED ALWAYS AS (won_sets - lost_sets) STORED,
    -- Sanal Sütun: Oyun Averajı
    game_diff INTEGER GENERATED ALWAYS AS (won_games - lost_games) STORED,

    CONSTRAINT unique_league_team UNIQUE (league_id, team_id),

    CONSTRAINT fk_scorboard_league
            FOREIGN KEY (league_id)
            REFERENCES tennisleague.league(id)
            ON DELETE CASCADE,

    CONSTRAINT fk_scorboard_team
            FOREIGN KEY (team_id)
            REFERENCES tennisleague.team(id)
            ON DELETE CASCADE,

    CONSTRAINT fk_scorboard_player
                FOREIGN KEY (player_id)
                REFERENCES tennisleague.player(id)
                ON DELETE CASCADE,

    -- Aynı ligde aynı oyuncu/takım iki kere olamaz
    CONSTRAINT unique_league_participant UNIQUE (league_id, player_id, team_id),

    -- Mantıksal kontrol: Sadece biri dolu olmalı
    CONSTRAINT chk_sb_participant CHECK (
        (player_id IS NOT NULL AND team_id IS NULL) OR
        (player_id IS NULL AND team_id IS NOT NULL)
    )
);

CREATE INDEX idx_sb_league_leaderboard ON tennisleague.score_board (
    league_id,
    score DESC,
    set_diff DESC,   -- Sanal kolon artık indekste!
    game_diff DESC   -- Eşitlik durumunda oyun averajına bakar
);