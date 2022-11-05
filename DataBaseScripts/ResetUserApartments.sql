\c moshack

DROP TABLE IF EXISTS user_apartments;

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
    total_price real,
    price_m2 real
);
