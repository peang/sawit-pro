-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

CREATE TABLE IF NOT EXISTS estates (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) UNIQUE,
    width INT NOT NULL CHECK (width >= 1 AND width <= 50000),
    length INT NOT NULL CHECK (length >= 1 AND length <= 50000),
    tree_count SMALLINT DEFAULT 0,
    min_tree_height SMALLINT,
    max_tree_height SMALLINT,
    median_tree_height SMALLINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_estates_uuid ON estates(uuid);

CREATE TABLE IF NOT EXISTS trees (
    id SERIAL PRIMARY KEY,
	uuid VARCHAR(36) UNIQUE,
    estate_id INTEGER REFERENCES estates(id),
    x INT NOT NULL CHECK (x >= 1), -- Assuming x and y are coordinates, which cannot be negative
    y INT NOT NULL CHECK (y >= 1), -- Assuming x and y are coordinates, which cannot be negative
    height SMALLINT NOT NULL CHECK (height >= 1 AND height <= 30),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_tree_location UNIQUE (estate_id, x, y), -- Ensure one tree per plot
    CONSTRAINT fk_estate_id FOREIGN KEY (estate_id) REFERENCES estates(id)
);

CREATE INDEX IF NOT EXISTS idx_trees_uuid ON trees(uuid);
CREATE INDEX IF NOT EXISTS idx_trees_x ON trees(x);
CREATE INDEX IF NOT EXISTS idx_trees_y ON trees(y);
CREATE INDEX IF NOT EXISTS idx_trees_estate_id ON trees(estate_id);
