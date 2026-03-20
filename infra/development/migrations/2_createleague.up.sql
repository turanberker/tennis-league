CREATE TABLE tennisleague.league (
    id  VARCHAR(100) PRIMARY KEY  DEFAULT gen_random_uuid(),
    name  VARCHAR(100) NOT NULL unique,

    --SINGLE,DOUBLE,TEAM
    format VARCHAR(20) NOT NULL,
    --MIX, MALE,FEMALE
    category VARCHAR(20) NOT NULL,
     --FIXTURE,DEFI
    process_type  VARCHAR(20) NOT NULL,
    -- Mevcut Durum: DRAFT, ACTIVE, COMPLETED
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',

    total_attendance INT8 NOT NULL DEFAULT 0,
    fixture_created_date TIMESTAMP,
    start_date TIMESTAMP,
    end_date TIMESTAMP
);