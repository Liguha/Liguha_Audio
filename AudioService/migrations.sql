-- создание таблицы пользаков
CREATE TABLE IF NOT EXISTS music.users(
                                          id serial primary key,
                                          login text not null,
                                          password text not null,
                                          role text not null
);

CREATE UNIQUE INDEX user_login ON music.users (login);
----------------

-- таблица на песни
CREATE TABLE IF NOT EXISTS music.songs(
                                          id serial primary key,
                                          name text not null,
                                          compositor text not null,
                                          author_id bigint not null
);

-- таблица альбомов
CREATE TABLE IF NOT EXISTS music.albums(
                                           id serial primary key,
                                           author_id bigint not null,
                                           name text not null,
                                           is_official bool default false
);

-- теги-песни
CREATE TABLE IF NOT EXISTS music.tags_songs(
                                               song_id bigint not null,
                                               tag_id bigint not null
);

CREATE TABLE IF NOT EXISTS music.tags(
                                         id serial primary key,
                                         name text not null
);

CREATE TABLE IF NOT EXISTS music.albums_songs(
                                                 album_id bigint not null,
                                                 song_id bigint not null
);