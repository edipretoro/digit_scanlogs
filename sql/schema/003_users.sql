-- +goose Up
create table if not exists users (
    id uuid primary key,
    uid integer not null unique,
    username text not null unique,
    fullname text not null unique,
    created_at timestamp not null default (CURRENT_TIMESTAMP),
    updated_at timestamp not null default (CURRENT_TIMESTAMP),
    deleted_at timestamp
);

-- +goose Down
drop table if exists users;