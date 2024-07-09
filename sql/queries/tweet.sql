-- name: CreateNewTweet :one
insert into tweets (id, content, media, public_id, user_id) values ($1, $2, $3, $4, $5) returning *;


-- name: FindATweet :one
select * from tweets where id=$1;

-- name: DeleteATweet :one
delete from tweets where id=$1 and user_id=$2 returning *;
