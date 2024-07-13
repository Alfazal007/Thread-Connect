-- name: CreateUser :one
insert into users (id, username, password, email, created_at, updated_at) values ($1, $2, $3, $4, $5, $6) returning *;

-- name: GetUserByName :one
select * from users where username=$1;

-- name: GetUniqueUser :one
select count(*) from users where username=$1 or email=$2;

-- name: GetUserById :one
select * from users where id=$1;

-- name: UpdateRefreshToken :one
update users set refresh_token=$1, updated_at=$2 where id=$3 returning *;

-- name: UpdatePassword :one
update users set password=$1,updated_at=$2 where id=$3 returning *;

