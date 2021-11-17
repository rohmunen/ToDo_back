CREATE TABLE users (
    id bigserial not null primary key,
    email varchar not null unique,
    name varchar not null,
    encrypted_password varchar not null
);

CREATE TABLE todos (
    id bigserial not null primary key,
    user_id integer,
    title varchar not null,
    body varchar not null,
    isDone boolean not null,
    isFavourite boolean not null,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE email_ver_hash (
    id bigserial not null primary key,
    email varchar not null,
    ver_hash varchar not null unique,
    expiration timestamp
);

CREATE TABLE todos_public (
    id bigserial not null primary key,
    todo_id integer,
    link_string varchar not null unique,
    FOREIGN KEY (todo_id) REFERENCES todos(id)
);