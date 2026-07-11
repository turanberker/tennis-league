-- 1. Önce eski kısıtlamayı kaldıralım
ALTER TABLE tennisleague.score_board DROP CONSTRAINT unique_league_participant;

-- 2. Bir ligde aynı oyuncu sadece bir kez var olabilir (team_id NULL iken)
CREATE UNIQUE INDEX uidx_league_player
    ON tennisleague.score_board (league_id, player_id)
    WHERE player_id IS NOT NULL;

-- 3. Bir ligde aynı takım sadece bir kez var olabilir (player_id NULL iken)
CREATE UNIQUE INDEX uidx_league_team
    ON tennisleague.score_board (league_id, team_id)
    WHERE team_id IS NOT NULL;