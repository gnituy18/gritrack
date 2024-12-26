CREATE TABLE users(
	username VARCHAR(32) NOT NULL PRIMARY KEY,
	email VARCHAR(320) NOT NULL UNIQUE,
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
	slug TEXT NOT NULL,
	display_name TEXT NOT NULL,
	position INTEGER NOT NULL,
	description TEXT NOT NULL DEFAULT "",
	public INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (username, slug)
);

CREATE TABLE tracker_entries(
	username,
	slug, 
	date TEXT NOT NULL,
	emoji TEXT NOT NULL DEFAULT "",
	content TEXT NOT NULL DEFAULT "",
	FOREIGN KEY(username, slug) REFERENCES trackers(username, slug),
	PRIMARY KEY (username, slug, date DESC)
);
