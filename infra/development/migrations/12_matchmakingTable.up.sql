-- 1. Matchmaking Talepleri Havuzu Tablosu
CREATE TABLE tennisleague.matchmaking_requests
(
    id             VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    league_id      VARCHAR(36) NULL,
    requester_id   VARCHAR(36) NOT NULL,
    requester_type VARCHAR(10) NOT NULL,
    target_date    DATE        NOT NULL,
    start_hour     INT         NOT NULL, -- Sadece başlangıç saati (Örn: 19, 20)
    status         VARCHAR(15) NOT NULL DEFAULT 'PENDING',
    created_at     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Veri tutarlılığı için kısıtlamalar
    CONSTRAINT chk_requester_type CHECK (requester_type IN ('SINGLE', 'DOUBLE','TEAM')),
    CONSTRAINT chk_status CHECK (status IN ('PENDING', 'MATCHED', 'CANCELLED')),
    CONSTRAINT chk_valid_hour CHECK (start_hour >= 0 AND start_hour <= 23)
);

-- Admin sorguları için güncellenmiş composite index
CREATE INDEX idx_mm_requests_search_v3
    ON tennisleague.matchmaking_requests (target_date, start_hour, status, league_id);

-- Mükerrer ilan engelleme indeksi
CREATE UNIQUE INDEX uq_mm_request_per_slot_v3
    ON tennisleague.matchmaking_requests (requester_id, target_date, start_hour, COALESCE(league_id, 'GLOBAL'))
    WHERE status = 'PENDING';