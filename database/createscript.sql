CREATE TABLE records(
    record_uuid uuid,
    account_uuid uuid,
    tags set<text>,
    words set<text>,
    record text,
    importance int,
    PRIMARY KEY ((record_uuid, account_uuid), importance)
) WITH CLUSTERING ORDER BY (importance DESC);

CREATE TABLE config(
    config_uuid uuid,
    name text,
    record text,
    PRIMARY KEY (config_uuid, name)
);