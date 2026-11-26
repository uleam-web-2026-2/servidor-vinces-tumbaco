# Notas de entorno Windows — Taller 13

- **Entorno:** Windows con PowerShell y Docker Desktop.
- **Codificacion de archivos:** Para evitar que > en Windows PowerShell 5 genere UTF-16, los archivos de texto y payloads temporales se escribieron con System.Text.UTF8Encoding(False) (UTF-8 sin BOM).
- **Permisos en esquema public de PostgreSQL (SQLSTATE 42501):** Al iniciar con pi_app sin privilegios de creacion en el esquema public, se otorgo GRANT CREATE ON SCHEMA public TO api_app; para permitir las migraciones/tablas iniciales.
- **Limite de peticiones en /login (429):** Ejecutar multiples inicios de sesion seguidos en ráfaga corta disparo 429 por limitador de tasa; espaciar los logins en las pruebas evito el bloqueo temporal.