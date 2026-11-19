# ADR: Gestión de Transacciones en el Backend de Go

## Contexto
El sistema requiere asegurar la atomicidad y consistencia en operaciones críticas (como el intercambio de libros y el registro de usuarios).

## Decisión
Se implementará el manejo explícito de transacciones (`database/sql` y GORM transactions) en la capa de servicios, permitiendo aislar las operaciones de lectura/escritura y aplicar un `Rollback` automático en caso de fallo.

## Consecuencias
- **Positivas:** Garantiza la integridad de los datos ante fallos concurrentes.
- **Negativas:** Incrementa ligeramente la complejidad en la firma de métodos al requerir el manejo de la transacción `*sql.Tx`.