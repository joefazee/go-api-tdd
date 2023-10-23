CREATE TABLE users(
  id bigserial PRIMARY KEY,
  name varchar(255) NOT NULL,
  email varchar(255) NOT NULL UNIQUE,
  password varchar(255) NOT NULL
);

CREATE TABLE posts(
      id bigserial PRIMARY KEY,
      title TEXT NOT NULL UNIQUE,
      body TEXT NOT NULL,
      user_id bigint not null,
      created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
      FOREIGN KEY (user_id) REFERENCES users (id)
);







