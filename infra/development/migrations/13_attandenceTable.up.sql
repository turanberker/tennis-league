BEGIN;

-- 1. ESKİ LİMİT TRIGGER'LARINI SİL (Artık 2 oyuncu sınırı yok)
DROP TRIGGER IF EXISTS trg_team_player_limit ON tennisleague.team_player;
DROP FUNCTION IF EXISTS tennisleague.check_team_player_limit();


-- 2. TABLO İSİMLERİNİ DEĞİŞTİR (Rename)
-- team -> attendance
ALTER TABLE tennisleague.team RENAME TO attendance;

-- team_player -> attendance_player
ALTER TABLE tennisleague.team_player RENAME TO attendance_player;


-- 3. KOLON İSİMLERİNİ VE KISITLAMALARI (CONSTRAINTS) GÜNCELLE
-- attendance_player tablosundaki 'team_id' kolonunu 'attendance_id' yapıyoruz
ALTER TABLE tennisleague.attendance_player RENAME COLUMN team_id TO attendance_id;

-- Eski foreign key isimlerini de yeni tablo yapısına uygun olarak değiştiriyoruz
ALTER TABLE tennisleague.attendance_player
    RENAME CONSTRAINT fk_team_player_team TO fk_attendance_player_attendance;

ALTER TABLE tennisleague.attendance_player
    RENAME CONSTRAINT fk_team_player TO fk_attendance_player_player;


-- 4. DUPLICATE PLAYER TRIGGER'INI YENİ YAPIYA GÖRE GÜNCELLE
-- Eski trigger'ı tamamen silelim
DROP TRIGGER IF EXISTS trg_dublicate_player_in_league ON tennisleague.attendance_player;
DROP FUNCTION IF EXISTS prevent_duplicate_player_in_league();

-- Yeni trigger fonksiyonunu 'attendance_player' ve 'attendance' tablolarına göre oluştur
CREATE OR REPLACE FUNCTION tennisleague.prevent_duplicate_player_in_league()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM tennisleague.attendance_player ap
        JOIN tennisleague.attendance a ON a.id = ap.attendance_id
        WHERE ap.player_id = NEW.player_id
          AND a.league_id = (
              SELECT league_id FROM tennisleague.attendance WHERE id = NEW.attendance_id
          )
    ) THEN
        RAISE EXCEPTION 'Oyuncu aynı ligde birden fazla defa katılamaz.';
END IF;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger'ı yeni isme sahip attendance_player tablosuna bağlıyoruz
CREATE TRIGGER trg_duplicate_player_in_league
    BEFORE INSERT ON tennisleague.attendance_player
    FOR EACH ROW
    EXECUTE FUNCTION tennisleague.prevent_duplicate_player_in_league();

COMMIT;