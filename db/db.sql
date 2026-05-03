CREATE TABLE media_sources (
    id          SERIAL PRIMARY KEY,
    full_name   TEXT    NOT NULL,
    url         TEXT    NOT NULL UNIQUE,
    trust_level TEXT    NOT NULL,
    is_blocked  BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE publications (
    id        SERIAL PRIMARY KEY,
    title     TEXT    NOT NULL,
    url       TEXT    NOT NULL,
    source_id INTEGER NOT NULL REFERENCES media_sources(id),
    status    TEXT    NOT NULL
);

CREATE TABLE incidents (
    id             SERIAL PRIMARY KEY,
    category       TEXT    NOT NULL,
    descript       TEXT    NOT NULL,
    severity       TEXT    NOT NULL,
    publication_id INTEGER REFERENCES publications(id)
);

INSERT INTO media_sources (full_name, url, trust_level, is_blocked) VALUES
    ('РБК',              'rbc.ru',       'high',   FALSE),
    ('Неизвестный блог', 'fakeblog.ru',  'low',    FALSE),
    ('ТАСС',             'tass.ru',      'high',   FALSE),
    ('Анонимный источник','anon-news.ru','medium',  TRUE);
