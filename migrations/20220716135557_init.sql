-- +goose Up

-- credentials
CREATE TABLE creds (
    id SERIAL PRIMARY KEY,
    username VARCHAR(30) NOT NULL,
    password VARCHAR(100) NOT NULL,
    role INTEGER NOT NULL
);

-- users
CREATE TABLE usr (
    id SERIAL PRIMARY KEY,
    username VARCHAR(30) NOT NULL,
    role INTEGER NOT NULL,
    nickname VARCHAR(30) NOT NULL,
    avatar_url VARCHAR(255) NOT NULL,
    xpoints INTEGER NOT NULL,
    registration_time TIMESTAMPTZ NOT NULL,
    last_daily_cards_update TIMESTAMPTZ NOT NULL
);

-- admins
CREATE TABLE admin (
    id SERIAL PRIMARY KEY,
    username VARCHAR(30) NOT NULL,
    role INTEGER NOT NULL
);

-- images
CREATE TABLE img (
    id SERIAL PRIMARY KEY,
    url VARCHAR(255) NOT NULL,
    type VARCHAR(30) NOt NULL
);

-- cards_static
CREATE TABLE card_static (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    goal VARCHAR(30) NOT NULL,
    type VARCHAR(30) NOT NULL,
    pool VARCHAR(30) NOT NULL,
    short_description TEXT,
    long_description TEXT,
    chain_name VARCHAR(255),
    chain_order INTEGER,
    max_progress INTEGER
);

CREATE TABLE award (
    id SERIAL PRIMARY KEY,
    card_static_id INTEGER REFERENCES card_static(id),
    xpoints INTEGER NOT NULL,
    prize VARCHAR(255) NOT NULL,
    prize_img_url VARCHAR(255) NOT NULL,
    opt FLOAT
);

INSERT INTO card_static
(title, created_at, goal, type, pool, short_description, long_description, chain_name, chain_order, max_progress)
VALUES
        -- ordinary daily
        -- 1
        ('title', NOW(), 'goal', 'ordinary', 'daily', 'short', 'long', NULL, NULL, NULL),
        -- 2
        ('title2', NOW(), 'goal2', 'ordinary', 'daily', 'short', 'long', NULL, NULL, NULL),
        -- 3
        ('title3', NOW(), 'goal3', 'ordinary', 'daily', 'short', 'long', NULL, NULL, NULL),
        -- ordinary const
        -- 4
        ('title5', NOW(), 'goal5', 'ordinary', 'const', 'short', 'long', 'chain0', 1, NULL),
        -- 5
        ('title4', NOW(), 'goal4', 'ordinary', 'const', 'short', 'long', 'chain0', 0, NULL),
        -- progress daily
        -- 6
        ('title6', NOW(), 'goal6', 'progress', 'daily', 'short', 'long', NULL, NULL, 10),
        -- 7
        ('title7', NOW(), 'goal7', 'progress', 'daily', 'short', 'long', NULL, NULL, 20),
        -- options daily
        -- 8
        ('title8', NOW(), 'goal8', 'options', 'daily', 'short', 'long', NULL, NULL, NULL),
        -- 9
        ('title9', NOW(), 'goal9', 'options', 'daily', 'short', 'long', NULL, NULL, NULL),
        -- options const
        -- 10
        ('title8', NOW(), 'goal8', 'options', 'const', 'short', 'long', 'chain0', 3, NULL),
        -- 11
        ('title9', NOW(), 'goal9', 'options', 'const', 'short', 'long', 'chain0', 2, NULL);

-- award
INSERT INTO award
(card_static_id, xpoints, prize, prize_img_url, opt)
VALUES
        -- 1
        (1, 1, 'prize', 'prize_img_url', NULL),
        -- 2
        (2, 2, 'prize', 'prize_img_url', NULL),
        -- 3
        (3, 3, 'prize', 'prize_img_url', NULL),
        -- 4
        (4, 4, 'prize', 'prize_img_url', NULL),
        -- 5
        (5, 5, 'prize', 'prize_img_url', NULL),
        -- 6
        (6, 6, 'prize', 'prize_img_url', NULL),
        -- 7
        (7, 7, 'prize', 'prize_img_url', NULL),

        -- 8 opt
        (8, 80, 'prize', 'prize_img_url', 100),
        (8, 90, 'prize', 'prize_img_url', 200),
        (8, 100, 'prize', 'prize_img_url', 300),
        -- 9 opt
        (9, 110, 'prize', 'prize_img_url', 10),
        (9, 210, 'prize', 'prize_img_url', 20),
        -- 10 opt
        (10, 110, 'prize', 'prize_img_url', 10),
        (10, 210, 'prize', 'prize_img_url', 20),
        -- 11 opt
        (11, 210, 'prize', 'prize_img_url', 20),
        (11, 110, 'prize', 'prize_img_url', 10);

-- +goose Down
DROP TABLE IF EXISTS usr, admin, img, creds, award, card_static;
