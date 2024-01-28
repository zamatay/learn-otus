-- +goose Up
-- +goose StatementBegin
CREATE TABLE calendar (
    id bigint NOT NULL,
    title text,
    date timestamp,
    date_interval int,
    description text,
    user_id int,
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists calendar;
-- +goose StatementEnd
