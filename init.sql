CREATE TABLE users(
	username VARCHAR(32) NOT NULL PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	timezone TEXT NOT NULL,
	public INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE user_sessions(
	username REFERENCES users,
	id VARCHAR(256),
	created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE trackers(
	username REFERENCES users(username),
	tracker_id TEXT NOT NULL,
	display_name TEXT NOT NULL,
	position INTEGER NOT NULL,
	description TEXT NOT NULL DEFAULT "",
	public INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (username, tracker_id)
);

CREATE TABLE tracker_entries(
	username,
	tracker_id,
	date TEXT NOT NULL,
	emoji TEXT NOT NULL DEFAULT "",
	content TEXT NOT NULL DEFAULT "",
	FOREIGN KEY(username, tracker_id)
	REFERENCES trackers(username, tracker_id)
	ON UPDATE CASCADE
	ON DELETE CASCADE,
	PRIMARY KEY (username, tracker_id, date DESC)
);
