CREATE TABLE IF NOT EXISTS fact(
    fact_id SERIAL PRIMARY KEY,
    related_fact_ids TEXT,
    related_facts TEXT,
    fact_data TEXT
);

CREATE TABLE IF NOT EXISTS config(
    id SERIAL PRIMARY KEY,
    key TEXT,
    value TEXT
);