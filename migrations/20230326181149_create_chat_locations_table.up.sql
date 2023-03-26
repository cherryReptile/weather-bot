create table chat_locations
(
    id           bigserial              not null primary key,
    username     varchar                not null unique,
    lng          float,
    lat          float,
    weather_stat varchar default '',
    country      varchar default '',
    city         varchar default '',
    chat_id      bigserial              not null unique,
    created_at   time without time zone not null,
    updated_at   time without time zone
)