-- name: CreateRole :one
INSERT INTO roles (name, priority, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListRoles :many
SELECT * FROM roles ORDER BY priority DESC;

-- name: AssignUserRole :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT (user_id, role_id) DO NOTHING;

-- name: GetUserRoles :many
SELECT r.id, r.name, r.priority, r.description
FROM roles r
JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: GetRoleByName :one
SELECT id, name, priority, description
FROM roles
WHERE name = $1;

-- name: AssignRoleToUser :exec
INSERT INTO user_roles (user_id, role_id, assigned_at)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, role_id) DO UPDATE SET assigned_at = EXCLUDED.assigned_at;

-- name: GetUsersByRole :many
SELECT u.id, u.name, u.email, u.created_at, u.updated_at
FROM users u
         JOIN user_roles ur ON u.id = ur.user_id
WHERE ur.role_id = $1;