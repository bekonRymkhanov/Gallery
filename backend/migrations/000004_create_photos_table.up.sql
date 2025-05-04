CREATE TABLE IF NOT EXISTS photos (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    description text,
    author text NOT NULL,
    category text NOT NULL,
    tags text,
    width integer NOT NULL,
    height integer NOT NULL,
    url text NOT NULL,
    thumbnail_url text NOT NULL,
    source text NOT NULL,
    download_count bigint DEFAULT 0,
    likes bigint DEFAULT 0,
    version integer NOT NULL DEFAULT 1
);

--indexes
CREATE INDEX IF NOT EXISTS photos_title_idx ON photos USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS photos_author_idx ON photos USING GIN (to_tsvector('simple', author));
CREATE INDEX IF NOT EXISTS photos_category_idx ON photos USING GIN (to_tsvector('simple', category));
CREATE INDEX IF NOT EXISTS photos_tags_idx ON photos USING GIN (to_tsvector('simple', tags));