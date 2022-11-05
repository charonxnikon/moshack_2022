CREATE DATABASE moshack;
\c moshack

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS user_apartments;

DROP TABLE IF EXISTS db_apartments;

CREATE TABLE users (
    id serial PRIMARY KEY,
    login varchar(25),
    password varchar(25)
);

CREATE TABLE user_apartments(
    id serial PRIMARY KEY,
    user_id serial, 
    address text,
    rooms text,
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
    longitude real,
    price real,
    price_sq_meter real
);

CREATE TABLE db_apartments(
    id serial PRIMARY KEY,
    user_id serial, 
    address text,
    rooms text,
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
    longitude real,
    price real,
    price_sq_meter real
);
