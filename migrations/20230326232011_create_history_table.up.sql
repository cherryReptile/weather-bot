create table history
(
    id         bigserial                   not null primary key,
    request    json                        not null,
    chat_id    bigserial                   not null,
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone
)