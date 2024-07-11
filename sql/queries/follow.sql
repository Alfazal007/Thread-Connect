-- name: CreateNewFollow :one
insert into follow (id, follower, following) values ($1, $2, $3) returning *;

-- name: UnfollowUser :one
delete from follow where follower=$1 and following=$2 returning *;

-- name: AlreadyFollowing :one
select * from follow where follower=$1 and following=$2;

