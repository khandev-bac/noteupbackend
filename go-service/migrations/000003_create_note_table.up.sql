CREATE TABLE IF NOT EXISTS notes(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    audio_url TEXT,
    audio_duration_seconds INT ,
    audio_file_size_mb INT,
    transcript TEXT,
    title TEXT,
    word_count INT,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','processing','completed')),
    search_vector TSVECTOR GENERATED ALWAYS AS (
          setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
          setweight(to_tsvector('english', coalesce(transcript, '')), 'B')
    ) STORED,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
