CREATE TABLE tennisleague.earned_points (
    id VARCHAR(100) PRIMARY KEY DEFAULT gen_random_uuid(),              -- Otomatik artan benzersiz kayıt ID'si
    player_id VARCHAR(100) NOT NULL,               -- Kullanıcı ID'si
    match_date TIMESTAMP NOT NULL,      -- Maçın gerçekleştiği tarih ve saat
    earned_point INT NOT NULL,          -- Kazanılan puan (Negatif değerler kaybı temsil eder)
    match_type VARCHAR(20) NOT NULL,     --  -- Maç Tipi: 'SINGLE' veya 'DOUBLE'

    CONSTRAINT earned_points_player FOREIGN KEY (player_id) REFERENCES tennisleague.player(id)
);

-- Sorgu performansını artırmak için sık kullanılan alanlara indeks ekleyelim
CREATE INDEX idx_player_id ON earned_points(player_id);
CREATE INDEX idx_match_date ON earned_points(match_date);