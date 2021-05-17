./opt/cassandra/bin/cqlsh

CREATE KEYSPACE martian
  WITH REPLICATION = { 
   'class' : 'SimpleStrategy', 
   'replication_factor' : 1 
  };

-- Get all of the entities and their associated records
CREATE TABLE entities_to_records(
    entity text,
    account_uuid uuid,
    record_uuid set<text>,
    PRIMARY KEY ((account_uuid, entity))
);

-- Get all of the words and their associated records
CREATE TABLE words_to_records(
    word text,
    account_uuid uuid,
    record_uuid set<text>,
    PRIMARY KEY ((account_uuid, word))
);

-- Get the record(s) back that we need to return
CREATE TABLE record(
    account_uuid uuid,
    record_uuid uuid,
    entities set<text>,
    words set<text>,
    record text,
    title text,
    importance int,
    PRIMARY KEY ((account_uuid, record_uuid), title)
);

CREATE TABLE config(
    config_uuid uuid,
    name text,
    record text,
    PRIMARY KEY (config_uuid, name)
);

-- https://medium.com/rahasak/cassandra-golang-client-14f50171846b