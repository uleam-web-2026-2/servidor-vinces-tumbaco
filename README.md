# Mesa de ayuda

> Reemplacen todo lo que está entre corchetes en su primer commit.
> Este README es la puerta de entrada del repositorio: en la semana 1 otra pareja debe poder levantar el servidor siguiendo solo lo que dice aquí, y desde la semana 4 es la base de la integración continua.

**Hilo Servidor · ULEAM · Período 2026-2**
Aplicaciones Web II (TDI-610) · Aplicación para el Servidor Web (IS-503)

## Integrantes

| Integrante | Usuario de GitHub | Paralelo |
| --- | --- | --- |
| Vinces Rivera Maria Jose | Majovinces | Web II B |
| Tumbaco Pillasagua Kenyis Yusley | Yusley21 | Web II B |

## El producto

El proyecto consiste en un Servidor Web de Mesa de Ayuda, diseñado para gestionar solicitudes y reportes de problemas de los usuarios. El sistema permitirá registrar, organizar y dar seguimiento a los tickets de soporte, facilitando la atención de las diferentes áreas encargadas de resolverlos.

## Cómo levantar el servidor

Primero se debe levantar PostgreSQL mediante Docker:
```bash
docker run --name pg -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:16
