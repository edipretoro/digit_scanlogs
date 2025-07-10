-- +goose Up
create table if not exists projects (
  id uuid primary key,
  name text not null,
  path text not null,
  description text,
  created_at timestamp not null default (CURRENT_TIMESTAMP),
  updated_at timestamp not null default (CURRENT_TIMESTAMP),
  delete_at timestamp
);

-- +goose Down
drop table if exists projects;