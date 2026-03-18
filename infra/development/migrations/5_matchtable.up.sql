CREATE TABLE tennisleague.match (
    id varchar(100) PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Maçın Nereden Geldiği (Hepsi Nullable)
    league_id VARCHAR(100) NULL,
    tournament_id VARCHAR(100) NULL,
    is_friendly BOOLEAN DEFAULT false,

    -- Maç Tipi: 'SINGLE' veya 'DOUBLE'
    match_type VARCHAR(20) NOT NULL,

    -- SINGLE maçlar için (Eğer match_type = 'SINGLE' ise buralar dolar)
    player_1_id VARCHAR(100) NULL,
    player_2_id VARCHAR(100) NULL,

    -- DOUBLE veya TEAM maçlar için (Eğer match_type = 'DOUBLE' ise buralar dolar)
    team_1_id VARCHAR(100) NULL,
    team_2_id VARCHAR(100) NULL,

    team_1_score INT,
    team_2_score INT,
    winner_id varchar(100),

    match_date TIMESTAMPTZ NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    approve_date  TIMESTAMPTZ NULL,
    -- FK'ları yine de koymanı öneririm, silme işlemlerinde sistem patlamaz
    CONSTRAINT fk_p1 FOREIGN KEY (player_1_id) REFERENCES tennisleague.player(id),
    CONSTRAINT fk_p2 FOREIGN KEY (player_2_id) REFERENCES tennisleague.player(id),
    CONSTRAINT fk_t1 FOREIGN KEY (team_1_id) REFERENCES tennisleague.team(id),
    CONSTRAINT fk_t2 FOREIGN KEY (team_2_id) REFERENCES tennisleague.team(id),

    CONSTRAINT chk_different_sides
        CHECK (
            -- SINGLE maç ise: Oyuncu kolonları dolu olmalı, takım kolonları boş olmalı
            (match_type = 'SINGLE' AND
             player_1_id IS NOT NULL AND player_2_id IS NOT NULL AND
             player_1_id <> player_2_id AND
             team_1_id IS NULL AND team_2_id IS NULL)

            OR

            -- DOUBLE veya TEAM maçı ise: Takım kolonları dolu olmalı, oyuncu kolonları boş olmalı
            (match_type IN ('DOUBLE', 'TEAM') AND
             team_1_id IS NOT NULL AND team_2_id IS NOT NULL AND
             team_1_id <> team_2_id AND
             player_1_id IS NULL AND player_2_id IS NULL)
        )
);



CREATE TABLE tennisleague.match_set (
    id varchar(100) PRIMARY KEY DEFAULT gen_random_uuid(),

    match_id varchar(100) NOT NULL,
    set_number INT NOT NULL,          -- 1,2,3...

    side_1_games INT ,
    side_2_games INT ,

    side_1_tie_break_score INT ,
    side_2_tie_break_score INT,

    CONSTRAINT fk_set_match
        FOREIGN KEY (match_id)
        REFERENCES tennisleague.match(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_fixture_set UNIQUE (match_id, set_number),

    CONSTRAINT chk_games_positive
        CHECK (side_1_games >= 0 AND side_2_games >= 0)
);