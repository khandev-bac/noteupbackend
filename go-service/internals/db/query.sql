-- name: Signup :one
INSERT INTO users (email, password,user_device)
VALUES ($1, $2 , $3)
RETURNING id, email, plan, user_device,created_at;

-- name: GoogleAuth :one
INSERT INTO users (email, google_id,picture,user_device)
VALUES ($1, $2,$3, $4)
ON CONFLICT (email)
DO UPDATE SET
    google_id = COALESCE(users.google_id, EXCLUDED.google_id),
    updated_at = now()
RETURNING id, email,plan,user_device,created_at;


-- name: GetUserByEmail :one
SELECT id,email,google_id,picture,is_active,plan,user_device,created_at
FROM users
WHERE email = $1;

-- name: GetUserByEmailLogin :one
SELECT *
FROM users
WHERE email = $1;


-- name: GetUserById :one
SELECT id,email,google_id,picture,is_active,plan,user_device,created_at
FROM users
WHERE id = $1;


-- name: CreateNotes :one
INSERT INTO notes(user_id,audio_url,status,audio_file_size_mb,audio_duration_seconds)
VALUES($1,$2,'processing',$3,$4)
RETURNING *;


-- name: AfterProcessingUpdateNotes :one
UPDATE notes
SET
    title = $1,
    transcript = $2,
    word_count = $3,
    status = 'completed',
    updated_at = NOW()
WHERE id = $4
RETURNING
    id,
    user_id,
    audio_url,
    audio_duration_seconds,
    audio_file_size_mb,
    transcript,
    title,
    word_count,
    status,
    search_vector,
    created_at,
    updated_at;



-- name: UpdateNoteWithNoteId :one
UPDATE notes
SET
    title = $1,
    transcript = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: GetAllUsersNotes :many
SELECT *
FROM notes
WHERE user_id = $1
ORDER BY created_at DESC;


-- name: GetNoteByNoteId :one
SELECT *
FROM notes
WHERE id = $1
    AND user_id = $2;

-- name: DeleteNoteById :exec
DELETE FROM notes
WHERE id = $1
    AND user_id = $2;


-- name: SearchNotes :many
SELECT *
FROM notes
WHERE user_id = $1
    AND search_vector @@ plainto_tsquery('english', $2)
ORDER BY ts_rank(
search_vector,
plainto_tsquery('english', $2)
) DESC;


-- name: GetUserCoinBalance :one
SELECT balance
FROM user_coin
WHERE user_id = $1;


-- name: DeductUserCoins :exec
UPDATE user_coin
SET balance = balance - $1,
    updated_at = now()
WHERE user_id = $2
  AND balance >= $1;

-- name: AddUserCoins :exec
UPDATE user_coin
SET balance = balance + $1,
      updated_at = now()
WHERE user_id = $2;

-- name: GetUserCoins :one
SELECT user_id, balance, updated_at
FROM user_coin
WHERE user_id = $1;


-- name: GetActiveCoinPacks :many
SELECT id, coin_value, coin_price, popular
FROM coin_packs
WHERE active = TRUE
ORDER BY popular DESC, coin_value ASC;




-- name: GetCoinPackById :one
SELECT *
FROM coin_packs
WHERE id = $1
  AND active = TRUE;



-- name: CreateUserCoins :exec
INSERT INTO user_coin (user_id, balance)
VALUES ($1, $2);



-- name: CreateCoinTransaction :exec
INSERT INTO coin_transactions (
    user_id,
    amount,
    reason
) VALUES (
    $1,
    $2,
    $3
);

-- name: GetUserCoinTransactions :many
SELECT
    id,
    user_id,
    amount,
    reason,
    created_at,
    updated_at
FROM coin_transactions
WHERE user_id = $1
ORDER BY created_at DESC;
