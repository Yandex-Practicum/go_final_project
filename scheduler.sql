CREATE TABLE scheduler(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment VARCHAR(256),
    repeat VARCHAR(128)
); 

CREATE INDEX scheduler_date ON scheduler (date);