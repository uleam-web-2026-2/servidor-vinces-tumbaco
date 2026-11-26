CREATE TABLE intercambios (
    id SERIAL PRIMARY KEY,
    estado VARCHAR(50) NOT NULL,
    usuario_ofrece_id INTEGER,
    usuario_recibe_id INTEGER
);