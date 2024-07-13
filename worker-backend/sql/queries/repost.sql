-- name: CountRepostActual :one
select count(*) from repost where tweet_id=$1;

-- name: AddNewRepost :one
INSERT INTO repost_count (tweet_id, count) VALUES ($1, $2) returning *;

-- name: UpdateRepost :one
update repost_count set count=$1 where tweet_id=$2 returning *;

-- name: DeleteRepost :one
delete from repost_count where tweet_id=$1 returning *;
