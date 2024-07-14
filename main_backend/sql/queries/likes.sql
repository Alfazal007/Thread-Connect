-- name: CreateNewLike :one
insert into likes (id, tweet_id, user_id) values ($1, $2, $3) returning *;

-- name: DeleteLike :one
delete from likes where user_id=$1 and tweet_id=$2 returning *;

-- name: GetLike :one
select * from likes where user_id=$1 and tweet_id=$2;
