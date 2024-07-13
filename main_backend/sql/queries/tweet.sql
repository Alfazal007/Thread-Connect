-- name: CreateNewTweet :one
insert into tweets (id, content, media, public_id, user_id, created_at) values ($1, $2, $3, $4, $5, $6) returning *;

-- name: CreateNewTweetNoMedia :one
insert into tweets (id, content, user_id, created_at) values ($1, $2, $3, $4) returning *;


-- name: FindATweet :one
select * from tweets where id=$1;

-- name: DeleteATweet :one
delete from tweets where id=$1 and user_id=$2 returning *;

-- name: GetPaginatedTweet :many
SELECT * FROM tweets ORDER BY created_at LIMIT $1 OFFSET $2;
