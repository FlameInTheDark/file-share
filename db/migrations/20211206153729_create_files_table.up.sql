create extension if not exists "uuid-ossp";

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