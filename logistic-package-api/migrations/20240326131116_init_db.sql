-- +goose Up
-- +goose StatementBegin
create table if not exists package
(
    id      bigserial
        constraint package_pk
            primary key,
    weight  bigint,
    title   varchar(50) not null,
    created timestamp   not null,
    updated timestamp,
    removed boolean
) partition by hash (id);

comment on table package is 'logistic package omp db';

comment on column package.weight is 'grams';

create table if not exists packages_1 partition of package
for values with (modulus 3, remainder 0);

create table if not exists packages_2 partition of package
for values with (modulus 3, remainder 1);

create table if not exists packages_3 partition of package
for values with (modulus 3, remainder 2);


alter table package
    owner to logistic;

create table if not exists package_events
(
    id         bigserial
        constraint package_events_pk
            primary key,
    package_id bigint  not null,
    type       integer not null,
    status     integer,
    payload    jsonb
);

comment on column package_events.payload is 'model.Package object';
alter table package
    owner to logistic;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS package;
DROP TABLE IF EXISTS package_events;
-- +goose StatementEnd
