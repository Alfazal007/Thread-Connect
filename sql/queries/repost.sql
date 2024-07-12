-- name: CreateNewRepost :one
insert into repost (id, tweet_id, user_id) values ($1, $2, $3) returning *;

-- name: DeleteRepost :one
delete from repost where id=$1 returning *;

-- name: GetRepost :one
select * from repost where user_id=$1 and tweet_id=$2;


-- name: GetRepostById :one
select * from repost where id=$1;
