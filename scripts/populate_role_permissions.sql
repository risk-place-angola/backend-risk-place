-- Populate role_permissions for fine-grained access control
-- Run this after seeding roles and permissions

BEGIN;

-- Admin: Full access to everything
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Citizen: Basic user permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON (
    (p.resource = 'report' AND p.action IN ('create', 'read', 'update')) OR
    (p.resource = 'user' AND p.action IN ('read', 'update')) OR
    (p.resource = 'risk_type' AND p.action = 'read')
)
WHERE r.name = 'citizen'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ERCE: Emergency responders - can verify and resolve reports
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON (
    (p.resource = 'report' AND p.action IN ('create', 'read', 'update', 'verify', 'resolve')) OR
    (p.resource = 'user' AND p.action IN ('read', 'update')) OR
    (p.resource = 'risk_type' AND p.action = 'read')
)
WHERE r.name = 'erce'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ERFCE: Fire emergency responders - same as ERCE
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON (
    (p.resource = 'report' AND p.action IN ('create', 'read', 'update', 'verify', 'resolve')) OR
    (p.resource = 'user' AND p.action IN ('read', 'update')) OR
    (p.resource = 'risk_type' AND p.action = 'read')
)
WHERE r.name = 'erfce'
ON CONFLICT (role_id, permission_id) DO NOTHING;

COMMIT;

-- Verify permissions
SELECT 
    r.name as role,
    p.resource,
    p.action,
    p.code
FROM roles r
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
ORDER BY r.name, p.resource, p.action;
