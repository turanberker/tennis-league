CREATE TABLE tennisleague.users (
    id VARCHAR(100) PRIMARY KEY  DEFAULT gen_random_uuid(),
    email  VARCHAR(100) NOT NULL UNIQUE,
    phone  VARCHAR(100),
    name  VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    password_hash TEXT NOT NULL,
    role  VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    approved BOOLEAN DEFAULT FALSE    
);

INSERT INTO tennisleague.users (email,phone,name,surname,password_hash,"role",created_at,approved) VALUES
	 ('turanberker@yahoo.com','','Mithat Berker','turanberker@yahoo.com','$2a$10$hKbtFfzTbmUgsH3yegmcROTaQ7HL9mPnQN75wnACJLvPnzifA82T6','ADMIN','2026-03-03 19:26:52.455',true);
