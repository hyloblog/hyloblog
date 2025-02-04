-- intialise table
DROP DATABASE IF EXISTS hyloblog_db;
CREATE DATABASE hyloblog_db;
-- connect to db
\c hyloblog_db

DROP USER IF EXISTS hyloblog_user;
CREATE USER hyloblog_user with encrypted password 'secret';

GRANT ALL PRIVILEGES ON DATABASE hyloblog_db TO hyloblog_user;

SET bytea_output = 'hex';

-- initialise schema
\i /docker-entrypoint-initdb.d/schema/schema.pg.sql
