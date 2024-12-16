-- +goose Up
-- +goose StatementBegin
alter table package
    drop column removed;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table package
    add removed boolean;
-- +goose StatementEnd
