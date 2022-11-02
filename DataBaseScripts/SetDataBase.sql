CREATE DATABASE gusev;
\c gusev
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS apartments;
CREATE TABLE users (
    id serial PRIMARY KEY,
    login varchar(25),
    password varchar(25)
);

CREATE TABLE apartments(
    id serial PRIMARY KEY,
    user_id serial, 
    address text,
    rooms smallint,
    type text,
    height smallint,
    material text,
    floor smallint,
    area real,
    kitchen real,
    balcony text,
    metro integer,
    condition text,
    latitude real,
    longitude real
);
