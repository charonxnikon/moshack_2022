\c gusev
DROP TABLE users;
CREATE TABLE users (
    id serial PRIMARY KEY,
    login varchar(25),
    password varchar(25)
);