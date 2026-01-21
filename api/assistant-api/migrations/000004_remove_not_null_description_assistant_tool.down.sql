UPDATE assistant_tools
SET description = ''
WHERE description IS NULL;

ALTER TABLE assistant_tools
ALTER COLUMN description SET NOT NULL;