CREATE TABLE "scheduler" (
	"id"	INTEGER,
	"date"	TEXT NOT NULL,
	"title"	TEXT NOT NULL,
	"comment"	TEXT,
	"repeat"	TEXT NOT NULL DEFAULT " ",
);

CREATE INDEX "scheduler_date" ON "scheduler" (
	"date"	DESC
);