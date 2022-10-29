CREATE DATABASE gusev;
\c gusev
CREATE TABLE users (
    id serial PRIMARY KEY,
    login varchar(25),
    password varchar(25)
);

CREATE TABLE apartments (
    id serial PRIMARY KEY,
    address varchar(100),
    rooms integer
);