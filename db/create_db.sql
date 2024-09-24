drop table if exists file_meta;
create table file_meta(
id serial primary key,
name text,
content_length int,
parts json,
status int
);
