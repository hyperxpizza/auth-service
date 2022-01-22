
create table users (
    id serial primary key,
    username varchar(100) not null,
    email varchar(100) not null,
    passwordHash text not null,
    created timestamp not null,
    updated timestamp not null
);