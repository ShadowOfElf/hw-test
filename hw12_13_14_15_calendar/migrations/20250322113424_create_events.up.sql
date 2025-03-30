CREATE TABLE events (
    id VARCHAR(255) PRIMARY KEY,
    title TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    duration INTERVAL NOT NULL,
    description TEXT,
    user_id INTEGER NOT NULL,
    notification_minute INTERVAL NOT NULL
);