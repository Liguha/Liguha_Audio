-- создание таблицы пользаков
CREATE TABLE IF NOT EXISTS music_users.users(
                                                id serial primary key,
                                                login text not null,
                                                password text not null,
                                                role text not null
);

CREATE UNIQUE INDEX user_login ON music_users.users (login)
----------------