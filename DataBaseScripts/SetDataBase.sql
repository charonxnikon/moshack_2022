CREATE DATABASE gusev;
\c gusev
DROP TABLE users;
DROP TABLE apartments;
CREATE TABLE users (
    id serial PRIMARY KEY,
    login varchar(25),
    password varchar(25)
);

CREATE TABLE apartments(
    id serial PRIMARY KEY,
    address varchar(100),
    rooms smallint,
    building_segment smallint,
    building_floors smallint,
    wall_material smallint,
    apartment_floor smallint,
    apartment_area real,
    kitchen_area real,
    balcony integer,
    metro_remoteness integer,
    condition integer,
    latitude real,
    longitude real
);
