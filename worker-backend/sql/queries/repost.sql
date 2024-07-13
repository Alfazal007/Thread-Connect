-- name: CountRepostActual :one
select count(*) from repost where tweet_id=$1;

-- name: UpdateRepost :one
INSERT INTO repost_count (tweet_id, count)
VALUES ($1, $2)
ON CONFLICT (tweet_id)
DO UPDATE SET count = $2 returning *;
