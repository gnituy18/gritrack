PRAGMA foreign_keys = ON;

CREATE TABLE user(
	username VARCHAR(64) NOT NULL PRIMARY KEY,
	email VARCHAR(320) NOT NULL UNIQUE,
	birthday TEXT NOT NULL,
	timezone TEXT NOT NULL
);

CREATE TABLE session(
	username REFERENCES user,
	id VARCHAR(256),
	created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tracker(
	username REFERENCES user,
	name VARCHAR(64) NOT NULL,
	position INTEGER NOT NULL,
	public INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (username, name)
);

CREATE TABLE day(
	username,
	tracker_name, 
	date TEXT NOT NULL,
	emoji TEXT NOT NULL DEFAULT "",
	content TEXT NOT NULL DEFAULT "",
	FOREIGN KEY(username, tracker_name) REFERENCES tracker(username, name),
	PRIMARY KEY (username, tracker_name, date DESC)
);
