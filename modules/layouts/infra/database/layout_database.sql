CREATE TABLE layouts
(
    uuid        varchar(36) primary key,
    name        varchar(60) not null,
    device_uuid varchar(36) not null,
    content     longblob    not null,
    created_at  datetime    not null,
    modified_at datetime    not null,

    constraint foreign key(device_uuid) references devices (uuid)
)

select * from layouts;

