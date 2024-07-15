-- name: GetFollowers :many
SELECT u.email FROM follow f JOIN users u ON f.follower = u.id WHERE f.following = $1;

-- name: GetUsername :one
select username from users where id=$1;
