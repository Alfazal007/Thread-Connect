-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    username text unique not null,
    password text not null,
    refresh_token text,
    email text unique not null,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CHECK (LENGTH(username) BETWEEN 5 AND 20),
    CHECK (LENGTH(password) >= 6)
);
CREATE TABLE tweets (
    id uuid PRIMARY KEY,
    content text,
    media text,
    public_id text,
    created_at TIMESTAMP NOT NULL,
    user_id uuid not null references users(id) on delete cascade,
    reply_tweet_id uuid references tweets(id),
    repost uuid references tweets(id) on delete cascade
);

CREATE TABLE repost_count (
  tweet_id uuid PRIMARY KEY references tweets(id) on delete cascade,
  count int default 0 
);


CREATE TABLE likes_count (
  tweet_id uuid PRIMARY KEY references tweets(id) on delete cascade,
  count int default 0 
);

CREATE TABLE repost (
  id uuid PRIMARY KEY,
  user_id uuid not null references users(id) on delete cascade,
  tweet_id uuid not null references tweets(id) on delete cascade,
  unique(user_id, tweet_id)
);

CREATE TABLE likes (
    id uuid PRIMARY KEY,
    user_id uuid not null references users(id) on delete cascade,
    tweet_id uuid not null references tweets(id) on delete cascade,
    unique(user_id, tweet_id)
);

CREATE TABLE follow (
    id uuid PRIMARY KEY,
    follower uuid not null references users(id) on delete cascade,
    following uuid not null references users(id) on delete cascade,
    unique(follower, following)
);

-- +goose Down
DROP TABLE follow;
DROP TABLE likes;
DROP TABLE repost;
DROP TABLE likes_count;
DROP TABLE repost_count;
DROP TABLE tweets;
DROP TABLE users;
