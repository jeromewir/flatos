-- name: GetFlats :many
SELECT * FROM flats;

-- name: UpsertFlat :one
INSERT INTO flats (
    id,
    external_id,
    name,
    address,
    city,
    price,
    zip_code,
    size,
    availability,
    status,
    created_at,
    updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
ON CONFLICT(external_id) DO NOTHING
RETURNING id;

