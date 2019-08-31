FROM postgres

WORKDIR /
RUN mkdir -p /docker-entrypoint-initdb.d
COPY createTable.sql  /docker-entrypoint-initdb.d/

RUN ls 
RUN cd docker-entrypoint-initdb.d/
WORKDIR /docker-entrypoint-initdb.d
RUN ls 
WORKDIR /
# RUN ["mkdir", "/docker-entrypoint-initdb.d"]
# ADD ./scripts/createTable.sql  /docker-entrypoint-initdb.d/
# CMD [ "psql -U postgres -d postgres -a -f /docker-entrypoint-initdb.d/createTable.sql" ]
# CMD ["psql -U postgres -d postgres -a -f /arex/createTable.sql"]