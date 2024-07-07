CREATE TABLE "scheduler" (
	 PRIMARY KEY AUTOINCREMENT "id",
	"date"	TEXT NOT NULL,
	"title"	TEXT NOT NULL,
	"comment"	TEXT,
	"repeat"	TEXT NOT NULL DEFAULT " ",
);

CREATE INDEX "scheduler_date" ON "scheduler" (
	"date"	DESC
);