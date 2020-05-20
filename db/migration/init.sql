-- asd
CREATE TABLE users
(
    id           serial primary key,
    email        text default null,
    tg_id        integer not null,
    payment_info text default null
);
