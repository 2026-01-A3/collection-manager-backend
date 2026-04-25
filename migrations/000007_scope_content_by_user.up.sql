-- Adiciona user_id em categorias, coleções e itens.
-- Etapas: ADD COLUMN nullable, backfill para o admin, SET NOT NULL.

ALTER TABLE categories  ADD COLUMN IF NOT EXISTS user_id INTEGER REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE collections ADD COLUMN IF NOT EXISTS user_id INTEGER REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE items       ADD COLUMN IF NOT EXISTS user_id INTEGER REFERENCES users (id) ON DELETE CASCADE;

UPDATE categories
   SET user_id = (SELECT id FROM users WHERE role = 'ADMIN' ORDER BY id LIMIT 1)
 WHERE user_id IS NULL;

UPDATE collections
   SET user_id = (SELECT id FROM users WHERE role = 'ADMIN' ORDER BY id LIMIT 1)
 WHERE user_id IS NULL;

UPDATE items
   SET user_id = (SELECT id FROM users WHERE role = 'ADMIN' ORDER BY id LIMIT 1)
 WHERE user_id IS NULL;

ALTER TABLE categories  ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE collections ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE items       ALTER COLUMN user_id SET NOT NULL;
