-- name: CreatePermission :one
INSERT INTO permissions (resource, action)
VALUES ($1, $2)
ON CONFLICT (resource, action) DO NOTHING
RETURNING *;

-- name: ListPermissions :many
SELECT * FROM permissions ORDER BY resource, action;

-- name: GetPermissionByCode :one
SELECT * FROM permissions WHERE code = $1;

-- name: AssignPermissionToRole :exec
INSERT INTO role_permissions (role_id, permission_id)
VALUES ($1, $2)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- name: GetRolePermissions :many
SELECT p.id, p.resource, p.action, p.code
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = $1;

-- name: GetUserPermissions :many
SELECT DISTINCT p.id, p.resource, p.action, p.code
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
JOIN user_roles ur ON rp.role_id = ur.role_id
WHERE ur.user_id = $1;

-- name: HasPermission :one
SELECT EXISTS(
    SELECT 1
    FROM permissions p
    JOIN role_permissions rp ON p.id = rp.permission_id
    JOIN user_roles ur ON rp.role_id = ur.role_id
    WHERE ur.user_id = $1 AND p.code = $2
) AS has_permission;

-- name: RemovePermissionFromRole :exec
DELETE FROM role_permissions
WHERE role_id = $1 AND permission_id = $2;
