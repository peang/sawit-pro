-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

-- Create table for estates
CREATE TABLE estates (
    id SERIAL PRIMARY KEY,
    width INTEGER NOT NULL CHECK (width >= 1 AND width <= 50000),
    length INTEGER NOT NULL CHECK (length >= 1 AND length <= 50000),
    tree_count INTEGER DEFAULT 0,
    min_height INTEGER,
    max_height INTEGER,
    median_height INTEGER
);

CREATE TABLE trees (
    id SERIAL PRIMARY KEY,
    estate_id INTEGER REFERENCES estates(id),
    x INTEGER NOT NULL CHECK (x >= 1), -- Assuming x and y are coordinates, which cannot be negative
    y INTEGER NOT NULL CHECK (y >= 1), -- Assuming x and y are coordinates, which cannot be negative
    height INTEGER NOT NULL CHECK (height >= 1 AND height <= 30), -- Height in meters
    CONSTRAINT unique_tree_location UNIQUE (estate_id, x, y), -- Ensure one tree per plot
    CONSTRAINT fk_estate_id FOREIGN KEY (estate_id) REFERENCES estates(id)
);

-- Create index on estate_id for faster lookup of trees by estate
CREATE INDEX idx_trees_estate_id ON trees(estate_id);

