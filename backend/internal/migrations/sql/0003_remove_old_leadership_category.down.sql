-- Revert: restore the old "Лидерство" category if needed.
-- Note: This is a one-way cleanup migration, but we provide a down migration for completeness.

INSERT INTO categories (name, slug, is_system, created_at)
VALUES ('Лидерство', 'leadership', TRUE, NOW())
ON CONFLICT (slug) DO NOTHING;
