PRAGMA foreign_keys = ON;

CREATE TABLE users(
	username VARCHAR(32) NOT NULL PRIMARY KEY,
	email VARCHAR(320) NOT NULL UNIQUE,
	timezone TEXT NOT NULL
);

CREATE TABLE user_sessions(
	username REFERENCES users,
	id VARCHAR(256),
	created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE trackers(
	username REFERENCES users,
	tracker_name TEXT NOT NULL,
	position INTEGER NOT NULL,
	description TEXT NOT NULL DEFAULT "",
	public INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (username, tracker_name)
);

CREATE TABLE tracker_entries(
	username,
	tracker_name, 
	date TEXT NOT NULL,
	emoji TEXT NOT NULL DEFAULT "",
	content TEXT NOT NULL DEFAULT "",
	FOREIGN KEY(username, tracker_name) REFERENCES trackers(username, tracker_name),
	PRIMARY KEY (username, tracker_name, date DESC)
);
