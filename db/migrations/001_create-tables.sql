CREATE TABLE eventlog (
                          eventid TEXT PRIMARY KEY,
                          name TEXT NOT NULL,
                          body BLOB,
                          metadata BLOB,
                          createdat TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);