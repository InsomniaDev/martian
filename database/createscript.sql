CREATE TABLE IF NOT EXISTS fact(
    id SERIAL PRIMARY KEY,
    name TEXT,
    related_fact_ids TEXT,
    related_facts TEXT,
    fact_data TEXT
);

CREATE TABLE IF NOT EXISTS config(
    id SERIAL PRIMARY KEY,
    key TEXT,
    value TEXT
);