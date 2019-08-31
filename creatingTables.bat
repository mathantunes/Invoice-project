docker exec -it postgres psql -U postgres

\c postgres

\ir /docker-entrypoint-initdb.d/createTable.sql