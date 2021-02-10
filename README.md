# martian
martian is the home ai system

## Connect to Postgres
- Create the table accounts
`CREATE TABLE accounts (user_id serial PRIMARY KEY,username VARCHAR ( 50 ) UNIQUE NOT NULL);`
- Open connection to the postgres database
`kubectl port-forward postgres-9c7b87574-jm2bf 5432`

## Login to postgres from the instance
- `psql -h 127.0.0.1 -U jarvis`