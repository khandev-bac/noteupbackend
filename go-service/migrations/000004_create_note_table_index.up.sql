CREATE INDEX IF NOT EXISTS notes_search_vector_idx
ON notes USING GIN (search_vector);

-- Title substring search (ILIKE %text%)
CREATE INDEX IF NOT EXISTS notes_title_trgm_idx
ON notes USING GIN (title gin_trgm_ops);

-- User lookup
CREATE INDEX IF NOT EXISTS notes_user_id_idx
ON notes (user_id);

-- Status filtering
CREATE INDEX IF NOT EXISTS notes_status_idx
ON notes (status);
