CREATE TABLE IF NOT EXISTS fact(
    id SERIAL PRIMARY KEY,
    name TEXT,
    related_fact_ids TEXT,
    related_facts TEXT,
    fact_data TEXT,
    importantance INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS facts_to_words(  
    id SERIAL PRIMARY KEY,
    fact_id INT,
    word_id INT,
    importance INT,
    importantance INT DEFAULT 0,
    UNIQUE (fact_id, word_id)
);

CREATE TABLE IF NOT EXISTS word(  
    id SERIAL PRIMARY KEY,
    word TEXT
);