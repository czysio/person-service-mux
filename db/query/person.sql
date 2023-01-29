-- name: CreatePerson :one
INSERT INTO people (
  first_name,
  surname,
  email,
  nickname,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetPersonById :one
SELECT * FROM people
WHERE id = $1 LIMIT 1;

-- name: GetPeople :many
SELECT * FROM people
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdatePerson :one
UPDATE people
set 
first_name = coalesce(sqlc.narg('first_name'), first_name), 
surname = coalesce(sqlc.narg('surname'), surname), 
email = coalesce(sqlc.narg('email'), email) ,
nickname = coalesce(sqlc.narg('nickname'), nickname), 
updated_at = coalesce(sqlc.narg('updated_at '), updated_at ) 
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeletePerson :exec
DELETE FROM people
WHERE id = $1;