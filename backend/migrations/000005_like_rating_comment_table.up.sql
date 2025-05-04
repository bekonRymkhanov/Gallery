CREATE TABLE IF NOT EXISTS comments (
    id bigserial PRIMARY KEY,
    photo_id bigint NOT NULL REFERENCES photos(id) ON DELETE CASCADE,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS comments_photo_id_idx ON comments(photo_id);
CREATE INDEX IF NOT EXISTS comments_user_id_idx ON comments(user_id);




CREATE TABLE IF NOT EXISTS ratings (
    id bigserial PRIMARY KEY,
    photo_id bigint NOT NULL REFERENCES photos(id) ON DELETE CASCADE,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    score integer NOT NULL CHECK (score >= 1 AND score <= 5),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1,
    UNIQUE (photo_id, user_id)
);

CREATE INDEX IF NOT EXISTS ratings_photo_id_idx ON ratings(photo_id);
CREATE INDEX IF NOT EXISTS ratings_user_id_idx ON ratings(user_id);




CREATE TABLE IF NOT EXISTS likes (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    photo_id bigint NOT NULL REFERENCES photos(id) ON DELETE CASCADE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1,
    UNIQUE (photo_id, user_id)
);
CREATE INDEX IF NOT EXISTS likes_photo_id_idx ON likes(photo_id);
CREATE INDEX IF NOT EXISTS likes_user_id_idx ON likes(user_id);