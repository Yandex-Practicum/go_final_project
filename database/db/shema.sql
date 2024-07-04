CREATE TABLE "scheduler" (
	"id"	INTEGER,
	"date"	TEXT NOT NULL,
	"title"	TEXT NOT NULL,
	"comment"	TEXT,
	"repeat"	TEXT NOT NULL DEFAULT " ",
	CHECK(length("repeat") <= 128)
	CHECK(length("title") > 0)
	PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE INDEX "scheduler_date" ON "scheduler" (
	"date"	DESC
);