DROP TABLE comments;
DROP TABLE ratings;
DROP TABLE likes;

DROP INDEX IF EXISTS comments_photo_id_idx;
DROP INDEX IF EXISTS comments_user_id_idx;

DROP INDEX IF EXISTS ratings_photo_id_idx;
DROP INDEX IF EXISTS ratings_user_id_idx;

DROP INDEX IF EXISTS likes_photo_id_idx;
DROP INDEX IF EXISTS likes_user_id_idx;