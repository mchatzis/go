-- name: CountOfTasksInState :one
SELECT COUNT(*) FROM tasks WHERE state = $1;

-- name: GetTotalValueOfDoneTasksByType :many
SELECT type, SUM(value) FROM tasks WHERE state='done' GROUP BY type ORDER BY type;

-- name: GetCountOfDoneTasksByType :many
SELECT type, COUNT(*) FROM tasks WHERE state='done' GROUP BY type ORDER BY type;

-- name: GetCountOfTasksByState :many
SELECT state, COUNT(*) FROM tasks GROUP BY state;

-- name: CreateTask :exec
INSERT INTO tasks (id, type, value, state, creationtime, lastupdatetime) VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateTaskState :exec
UPDATE tasks SET state = $1 WHERE id = $2;