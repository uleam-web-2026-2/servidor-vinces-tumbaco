-- name: BuscarIntercambio :one
SELECT id, estado
FROM intercambios
WHERE id = $1;

-- name: GuardarIntercambio :execrows
UPDATE intercambios
SET estado = $2
WHERE id = $1;
