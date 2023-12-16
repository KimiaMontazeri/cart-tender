CREATE TABLE IF NOT EXISTS users
(
    username varchar(250) NOT NULL,
    password varchar(250) NOT NULL,
    PRIMARY KEY
        (
         username
            )
);

CREATE TYPE cart_state AS ENUM ('COMPLETED', 'PENDING');

CREATE TABLE IF NOT EXISTS carts
(
    id         serial       NOT NULL,
    username   varchar(250) NOT NULL,
    created_at timestamp    NOT NULL DEFAULT NOW(),
    updated_at timestamp             DEFAULT NULL,
    data       JSONB,
    state      cart_state            DEFAULT 'PENDING',
    PRIMARY KEY
        (
         id
            )
);