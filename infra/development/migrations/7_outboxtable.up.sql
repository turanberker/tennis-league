CREATE TABLE tennisleague.outbox_events (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    aggregate_type  varchar(100) NOT NULL,   -- örn: 'match'
    aggregate_id    varchar(100) NOT NULL,   -- örn: matchId
    event_type      varchar(100) NOT NULL,   -- örn: 'MatchApproved'

    payload         jsonb NOT NULL,

    status          varchar(20) NOT NULL DEFAULT 'PENDING',
    retry_count     int NOT NULL DEFAULT 0,

    created_at      timestamp NOT NULL DEFAULT now(),
    processed_at    timestamp NULL
);

CREATE INDEX idx_outbox_pending ON tennisleague.outbox_events (status, created_at);