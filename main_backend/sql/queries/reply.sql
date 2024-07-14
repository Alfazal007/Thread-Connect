-- name: CreateNewReply :one
insert into tweets (id, content, user_id, created_at, reply_tweet_id) values ($1, $2, $3, $4, $5) returning *;

-- name: DeleteReply :one
delete from tweets where id=$1 and user_id=$2 returning *;
