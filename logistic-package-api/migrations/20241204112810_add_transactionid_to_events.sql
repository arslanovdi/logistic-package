-- +goose Up
-- +goose StatementBegin
alter table package_events
    add traceid CHAR(32);

comment on column package_events.traceid is 'OpenTelemetry root traceID';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table package_events
drop column traceid;
-- +goose StatementEnd
