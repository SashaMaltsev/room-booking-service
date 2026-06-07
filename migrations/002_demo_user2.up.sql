BEGIN;

INSERT INTO users (id, email, role)
VALUES ('00000000-0000-0000-0000-000000000003', 'dummy-user-2@example.com', 'user')
ON CONFLICT (id) DO NOTHING;

COMMIT;
