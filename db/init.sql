CREATE TABLE "user" (
    id bigint GENERATED ALWAYS AS IDENTITY,
    username text UNIQUE,
    isOnline bool DEFAULT false,
    password_hash text,
    bio text DEFAULT '',
    createdAt timestamp,
    updatedAt timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE chat (
    id bigint GENERATED ALWAYS AS IDENTITY,
    users bigint[],
    createdAt timestamp DEFAULT CURRENT_TIMESTAMP,
    updatedAt timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE message (
    id bigint GENERATED ALWAYS AS IDENTITY,
    userId bigint,
    chatId bigint,
    message text,
    isRead bool DEFAULT false,
    createdAt timestamp DEFAULT CURRENT_TIMESTAMP
);
