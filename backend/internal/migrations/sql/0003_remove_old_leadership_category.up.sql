-- Remove the old duplicate "Лидерство" category (slug='leadership')
-- that was created in migration 0001. We only want "Курс Лидерство" (slug='leadership-course').

DELETE FROM categories WHERE slug = 'leadership';
