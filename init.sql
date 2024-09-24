CREATE TABLE user(
	username VARCHAR(64) NOT NULL PRIMARY KEY,
	email VARCHAR(320) NOT NULL UNIQUE,
	birthday TEXT NOT NULL,
	public INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE session(
	username REFFERENCES user,
	id VARCHAR(256),
	created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tracker(
	username REFFERENCES user,
	name VARCHAR(64) NOT NULL,
	PRIMARY KEY (username, name)
);

CREATE TABLE day(
	username,
	tracker_name, 
	date TEXT NOT NULL,
	content TEXT,
	FOREIGN KEY(username, tracker_name) REFERENCES tracker(username, name),
	PRIMARY KEY (username, tracker_name, date DESC)
);
