CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);
CREATE TABLE IF NOT EXISTS users_permissions (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);
-- Add the two permissions to the table.
INSERT INTO permissions (code)
VALUES
    ('movies:read'),
    ('movies:write');




-- Add the admin user to the table.
-- Note: The password hash is for the password 'adminbek'.
INSERT INTO users (name, activated, is_admin, email, password_hash)
VALUES (
    'admin',
    true,
    true,
    'admin@example.com',
    '\x24326124313224586672434c73545676777a4b79424149664b7251412e6533674c52324a4d444d4f57712f346a345853456f4f70687556763350696d'
)
RETURNING id;

INSERT INTO users_permissions (user_id, permission_id)
VALUES
    (1, 1),
    (1, 2);