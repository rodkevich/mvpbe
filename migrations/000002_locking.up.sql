BEGIN;

CREATE TABLE Lock (
    lock_id VARCHAR(100) PRIMARY KEY,
    expires TIMESTAMP NOT NULL
);

END;
