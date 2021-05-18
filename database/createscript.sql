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
    name text,
    record text,
    PRIMARY KEY (name)
);

INSERT INTO config(name, record)
VALUES ('commonWords', 'the,at,there,some,my,of,be,use,her,than,and,this,an,would,first,a,have,each,make,water,to,from,which,like,been,in,or,she,him,call,is,one,do,into,who,you,had,how,time,oil,that,by,their,has,its,it,word,if,look,now,he,but,will,two,find,was,not,up,more,long,for,what,other,write,down,on,all,about,go,day,are,were,out,see,did,as,we,many,number,get,with,when,then,no,come,his,your,them,way,made,they,can,these,could,may,I,said,so,people,part');

-- https://medium.com/rahasak/cassandra-golang-client-14f50171846b