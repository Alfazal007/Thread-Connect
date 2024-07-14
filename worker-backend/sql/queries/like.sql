-- name: CountLikesActual :one
select count(*) from likes where tweet_id=$1;

-- name: AddNewLike :one
INSERT INTO likes_count (tweet_id, count) VALUES ($1, $2) returning *;

-- name: UpdateLike :one
update likes_count set count=$1 where tweet_id=$2 returning *;

-- name: DeleteLike :one
delete from likes_count where tweet_id=$1 returning *;
