-- --------------------- AUTHENTICATION ------------------------------
-- name: GetUserByEmail :one
SELECT user_id, user_name, email, phone, pass, pss_valid, otp, user_valid, otp_exp, role 
FROM news.users 
WHERE email = $1;

-- name: GetUserByUserName :one
SELECT user_id, user_name, email, phone, pass, pss_valid, otp, user_valid, otp_exp, role 
FROM news.users 
WHERE user_name = $1;

-- name: GetUserByPhone :one
SELECT user_id, user_name, email, phone, pass, pss_valid, otp, user_valid, otp_exp, role 
FROM news.users 
WHERE phone = $1;

-- name: GetUserById :one
SELECT user_id, user_name, email, phone, pass, pss_valid, otp, user_valid, otp_exp, role 
FROM news.users 
WHERE user_id = $1;

-- name: GetUserByLogin :one
SELECT user_id, user_name, email, phone, pass, pss_valid, otp, user_valid, otp_exp, role 
FROM news.users 
WHERE user_name = $1 OR email = $1 OR phone = $1;

-- name: CreateUser :exec
INSERT INTO news.users(user_name, email, phone, pass, otp, user_valid, otp_exp, role) 
VALUES($1, $2, $3, $4, $5, false, $6, $7);

-- name: SendNewOtp :exec
UPDATE news.users
SET otp = $1, otp_exp = $2 WHERE email = $3;

-- name: ActivateUser :exec
UPDATE news.users
SET user_valid = true WHERE email = $1;

-- name: UpdatePassword :exec
UPDATE news.users 
SET pass = $1, pss_valid = $2 
WHERE email = $3;

-- name: UpdateUser :exec
UPDATE news.users 
SET user_name = $1, email = $2, phone = $3, role = $4
WHERE user_id = $5;

-- name: GetAllUsers :many
SELECT user_id, user_name, email 
FROM news.users 
ORDER BY user_id;

-- --------------------- Category DATA ------------------------------
-- name: CreateCategory :one
INSERT INTO news.categories (name, slug)
VALUES ($1, $2) RETURNING *;

-- name: GetCategory :one
SELECT * FROM news.categories WHERE id = $1 LIMIT 1;

-- name: UpdateCategory :one
UPDATE news.categories 
SET name = $2, slug = $3 WHERE id = $1 
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM news.categories WHERE id = $1;

-- name: CheckCategoryExists :one
SELECT EXISTS(
    SELECT 1 FROM news.categories 
    WHERE id = $1
);

-- name: CheckCategorySlugExists :one
SELECT EXISTS(
    SELECT 1 FROM news.categories 
    WHERE slug = $1 OR name = $2
);

-- name: CheckCategoryInUse :one
SELECT EXISTs(
    SELECT 1 FROM news.articles
    WHERE category_id = $1
);

-- name: GetAllCategory :many
SELECT * from news.categories ORDER BY name;

-- ----------------- Comments Data ---------------------------------
-- name: CreateCommentWithDefaults :one
INSERT INTO news.comments (article_id, user_name, user_email, content, is_approved)
VALUES ($1, $2, $3, $4, true)
RETURNING *;

-- name: GetApprovedCommentsByArticle :many
SELECT * FROM news.comments 
WHERE article_id = $1 AND is_approved = true 
ORDER BY created_at DESC;

-- name: GetPendingComments :many
SELECT * FROM news.comments 
WHERE is_approved = false 
ORDER BY created_at DESC;

-- name: ApproveComment :one
UPDATE news.comments 
SET is_approved = true WHERE id = $1 
RETURNING *;

-- name: DisableComment :one
UPDATE news.comments 
SET is_approved = false WHERE id = $1 
RETURNING *;


-- name: ApprovalDueListOfComments :many
SELECT c.id AS comment_id, a.title , a."content" AS news,
c.user_name, c.user_email , c."content" AS comment 
FROM news."comments" c 
JOIN news.articles a ON a.id = c.article_id 
WHERE c.is_approved = false;

-- ------------------- Article Service ------------------------------

-- name: GetArticleDetails :one
SELECT * FROM news.articles
WHERE id = $1;

-- name: GetApprovedEngArticles :many
SELECT * FROM news.articles
WHERE status = 'published' AND category_id != 7 ORDER BY published_at DESC
OFFSET $1 LIMIT $2;

-- name: GetApprovedBanArticles :many
SELECT * FROM news.articles
WHERE status = 'published' AND category_id= 7 ORDER BY published_at DESC
OFFSET $1 LIMIT $2;

-- name: GetTodaysBnNews :many
SELECT *
FROM news.articles
WHERE status = 'published'
    AND category_id = 7 
    AND published_at >= $1::date
    AND published_at <  $1::date + INTERVAL '1 day'
ORDER BY published_at DESC;

-- name: GetTodaysEnNews :many
SELECT *
FROM news.articles
WHERE status = 'published'
    AND category_id != 7 
    AND published_at >= $1::date
    AND published_at <  $1::date + INTERVAL '1 day'
ORDER BY published_at DESC;

-- name: ReadArticleCount :exec
UPDATE news.articles
SET views_count = views_count + 1
WHERE id = $1;

-- name: GetUnApprovedArticleList :many
SELECT * FROM news.articles
WHERE status = 'draft' order BY published_at DESC;
-- OFFSET $1 LIMIT $2;

-- name: ApproveArticle :exec
UPDATE news.articles
SET status = 'published' WHERE id =$1;

-- name: DraftArticle :exec
UPDATE news.articles
SET status = 'draft' WHERE id =$1;