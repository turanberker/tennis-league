CREATE TABLE tennisleague.teams (
     id VARCHAR(100) PRIMARY KEY ,
     league_id VARCHAR(100) NOT NULL,
     name VARCHAR(100) NOT NULL,
      CONSTRAINT fk_team_league
             FOREIGN KEY (league_id)
             REFERENCES tennisleague.leagues(id)
);



CREATE TABLE tennisleague.team_players (
    team_id VARCHAR(100) NOT NULL,
    player_id VARCHAR(100) NOT NULL,

    PRIMARY KEY (team_id, player_id),

    CONSTRAINT fk_team_players_team FOREIGN KEY (team_id) REFERENCES tennisleague.teams(id) ON DELETE CASCADE,
    CONSTRAINT fk_team_players FOREIGN KEY (player_id) REFERENCES tennisleague.players(id)
);

CREATE OR REPLACE FUNCTION tennisleague.check_team_player_limit()
RETURNS TRIGGER AS $$
BEGIN
    IF (
        SELECT COUNT(*)
        FROM tennisleague.team_players
        WHERE team_id = NEW.team_id
    ) >= 2 THEN
        RAISE EXCEPTION 'Bir takım en fazla 2 oyuncudan oluşabilir';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_team_player_limit
BEFORE INSERT ON tennisleague.team_players
FOR EACH ROW
EXECUTE FUNCTION tennisleague.check_team_player_limit();


CREATE OR REPLACE FUNCTION prevent_duplicate_player_in_league()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM team_players tp
        JOIN teams t ON t.id = tp.team_id
        WHERE tp.player_id = NEW.player_id
          AND t.league_id = (
              SELECT league_id FROM teams WHERE id = NEW.team_id
          )
    ) THEN
        RAISE EXCEPTION 'Oyuncu aynı ligde birden fazla takımda olamaz';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_dublicate_player_in_league
BEFORE INSERT ON tennisleague.team_players
FOR EACH ROW
EXECUTE FUNCTION tennisleague.prevent_duplicate_player_in_league();