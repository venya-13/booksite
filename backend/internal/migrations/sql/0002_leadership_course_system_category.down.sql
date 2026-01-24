-- Revert: keep the categories table but remove the specific system row.

DELETE FROM categories WHERE slug = 'leadership-course';

