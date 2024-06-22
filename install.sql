CREATE TABLE  scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT id,
    date VARCHAR,
    title VARCHAR,
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX idx_date ON scheduler (date);