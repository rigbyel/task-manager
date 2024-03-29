CREATE TABLE IF NOT EXISTS quest (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    cost INT NOT NULL
);

CREATE TABLE IF NOT EXISTS user(
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    balance INTEGER
);

CREATE TABLE IF NOT EXISTS quest_stat(
    id INTEGER PRIMARY KEY,
    userID INTEGER NOT NULL,
    questID INTEGER NOT NULL,
    FOREIGN KEY (userID) REFERENCES user(id) ON DELETE CASCADE
);