-- Ensure the system category for the leadership course always exists and has the expected name/slug.
-- Canonical identifier: slug = 'leadership-course'

INSERT INTO categories (name, slug, is_system, created_at)
VALUES ('Курс Лидерство', 'leadership-course', TRUE, NOW())
ON CONFLICT (slug) DO UPDATE
SET name = EXCLUDED.name,
    is_system = TRUE;

