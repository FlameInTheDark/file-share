create type file_status as enum (
    'created',
    'uploaded'
);

create table files
(
    id         serial primary key,
    status     file_status not null        default 'created',
    file_id    uuid        not null unique default uuid_generate_v4(),
    file_name  varchar     not null        default 'file',
    downloads  bigint      not null        default 0,
    created_at timestamp   not null        default now()
);

-- name: Find :one
select *
from files
where file_id = $1;

-- name: Create :one
insert into files (file_name)
values ($1) returning *;

-- name: Uploaded :one
update files
set status = 'uploaded'
where file_id = $1 returning *;

-- name: Increase :exec
update files
set downloads = downloads + 1
where id = $1;

-- name: Stats :one
select downloads from files where file_id = $1;