drop table if exists file_meta;
create table file_meta(
    id serial primary key,
    name text NOT NULL,
    content_length int NOT NULL,
    parts json NOT NULL,
    status int NOT NULL,
    created_at timestamp NOT NULL default now()
);
