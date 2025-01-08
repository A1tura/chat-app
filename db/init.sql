CREATE TABLE "user" (
    id bigint GENERATED ALWAYS AS IDENTITY,
    username text UNIQUE,
    password_hash text,
    bio text DEFAULT '',
    createdAt timestamp,
    updatedAt timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE chat (
    id bigint GENERATED ALWAYS AS IDENTITY,
    users bigint[],
    createdAt timestamp,
    updatedAt timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE message (
    id bigint GENERATED ALWAYS AS IDENTITY,
    userId bigint,
    chatId bigint,
    message text,
    createdAt timestamp
);
