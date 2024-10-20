CREATE TABLE IF NOT EXISTS users
(
    uuid        varchar(36) primary key,
    name        varchar(60) not null,
    email       varchar(60) not null,
    password    text        not null,
    user_type   int         not null,
    created_at  datetime    not null,
    modified_at datetime    not null
);

ALTER TABLE users ADD COLUMN nickname VARCHAR(60) NOT NULL;

ALTER TABLE users MODIFY COLUMN password LONGBLOB NOT NULL;

select * from users