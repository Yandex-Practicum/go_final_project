CREATE TABLE scheduler (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  date VARCHAR(8) NOT NULL,
  title VARCHAR(256) NOT NULL,
  comment TEXT,
  repeat VARCHAR(128)
);

CREATE INDEX date_index ON scheduler (date);