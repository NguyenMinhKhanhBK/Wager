FROM mysql
COPY sql_migration/*.sql /docker-entrypoint-initdb.d/
