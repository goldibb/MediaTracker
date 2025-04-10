CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT,
    isbn TEXT,
    publication_year int,
    description TEXT,
    image_url TEXT,
    read boolean DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    rating_id INTEGER,
    watched BOOLEAN DEFAULT FALSE,
    theater_release_date TEXT,
    home_release_date TEXT,
    image_url TEXT,
    description TEXT,
    user_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS shows (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    rating_id INTEGER,
    episodes_aired INTEGER DEFAULT 0,
    pause_status BOOLEAN DEFAULT FALSE,
    image_url TEXT,
    description TEXT,
    user_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS episodes (
    id SERIAL PRIMARY KEY,
    show_id INTEGER NOT NULL,
    episode_number INTEGER NOT NULL,
    watched BOOLEAN DEFAULT FALSE,
    skipped BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (show_id) REFERENCES shows(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS ratings (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    priority INTEGER DEFAULT 0
);
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    email TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);