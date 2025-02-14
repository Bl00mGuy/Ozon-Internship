CREATE TABLE IF NOT EXISTS urls (
                                    id SERIAL PRIMARY KEY,
                                    original_url TEXT NOT NULL UNIQUE,
                                    short_url TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_original_url ON urls (original_url);
CREATE INDEX IF NOT EXISTS idx_short_url ON urls (short_url);