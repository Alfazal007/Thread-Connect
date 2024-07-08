-- name: CreateNewTweet :one
insert into tweets (id, content, media, public_id, user_id) values ($1, $2, $3, $4, $5) returning *;

