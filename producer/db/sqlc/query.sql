-- name: GetTasks :many
SELECT * FROM tasks;

-- name: GetTaskById :one
SELECT * FROM tasks WHERE id = $1;

-- name: CreateTask :exec
INSERT INTO tasks (id, type, value, state, creationtime, lastupdatetime) VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateTask :exec
UPDATE tasks SET id = $1, type = $2, value = $3, state = $4, lastupdatetime = $5 WHERE id = $6;