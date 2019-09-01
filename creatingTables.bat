REM This might be useful if the PostgreSQL initialiazation configuration does not kick in

docker exec -it postgres psql -U postgres

\c postgres

\ir /docker-entrypoint-initdb.d/createTable.sql