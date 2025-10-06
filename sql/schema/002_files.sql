-- +goose Up
create table if not exists files (
    id uuid primary key,
    project_id uuid not null references projects(id) on delete cascade,
    user_id uuid not null references users(id) on delete cascade,
    name text not null,
    path text not null,
    size bigint not null,
    mode text not null,
    modtime timestamp not null,
    sha512 text not null,
    description text,
    created_at timestamp not null default (CURRENT_TIMESTAMP),
    updated_at timestamp not null default (CURRENT_TIMESTAMP),
    deleted_at timestamp
);

-- +goose Down
drop table if exists files;