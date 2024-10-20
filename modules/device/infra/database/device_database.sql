CREATE TABLE devices
(
    uuid        varchar(36) primary key,
    name        varchar(60) not null,
    height      numeric     not null,
    width       numeric     not null,
    orientation numeric     not null,
    created_at  datetime    not null,
    modified_at datetime    not null
)

select * from devices