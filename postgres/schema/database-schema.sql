drop table if exists users;
create table users (
    id serial primary key,
    username varchar(100) unique not null,
    passwordHash text not null,
    created timestamp not null,
    updated timestamp not null,
    relatedUsersServiceID int not null
);
