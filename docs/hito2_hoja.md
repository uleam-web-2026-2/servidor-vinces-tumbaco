# Hoja de Seguimiento - Hito 2

**Proyecto:** KEMA (Mesa de Ayuda / Intercambios de Libros Usados)
**Integrantes:** 
- María José Vinces Rivera (`Majovinces`)
- Kenyis Yusley Tumbaco Pillasagua (`Yusley21`)

---

## 1 · Usuarios de ejemplo

| Usuario | Rol | Contraseña de ejemplo |
| :--- | :--- | :--- |
| `majo_user` | usuario | `user123` |
| `kenyis_admin` | administrador | `admin123` |

*Son usuarios de prueba con credenciales asignadas para la demostración en vivo.*

---

## 2 · Las seis peticiones

| # | Petición | Con qué credencial | Debe responder | Captura de respaldo |
| :---: | :--- | :--- | :---: | :--- |
| **0** | `POST /intercambios` (Preparación) | Token Admin | `201` / `200` | `0_0_preparar.png` |
| **1** | `GET /intercambios` | Ninguna | `401` | `1_1_sin_token.png` |
| **2** | `POST /login` | Ninguna | `200` y el token | `2_2_login.png` |
| **3** | `GET /intercambios` | Token de `majo_user` | `200` | `3_3_con_token.png` |
| **4** | `PATCH /intercambios/1` (Pendiente → Aceptado) | Token de `majo_user` (usuario) | `403` | `4_4_rol_sin_permiso.png` |
| **5** | `PATCH /intercambios/1` (Pendiente → Aceptado) | Token de `kenyis_admin` (admin) | `200` | `5_5_rol_correcto.png` |
| **6** | `PATCH /intercambios/1` (Pendiente → Entregado) | Token de `kenyis_admin` (admin) | `422` | `6_6_transicion_prohibida.png` |

---

## 3 · Matriz ruta × rol

| Ruta | Sin sesión | Usuario | Administrador |
| :--- | :---: | :---: | :---: |
| `POST /login` | Sí | Sí | Sí |
| `GET /salud` | Sí | Sí | Sí |
| `GET /intercambios` | No (401) | Sí | Sí |
| `GET /intercambios/{id}` | No (401) | Sí | Sí |
| `POST /intercambios` | No (401) | No (403) | Sí |
| `PUT /intercambios/{id}` | No (401) | No (403) | Sí |
| `PATCH /intercambios/{id}` | No (401) | No (403) | Sí |
| `DELETE /intercambios/{id}` | No (401) | No (403) | Sí |

---

## 4 · Recuadro de cumplimiento

- ☑ Los structs de todas las entidades están en el código y sus tablas existen.
- ☑ El validador termina sin errores ni avisos.
- ☑ Las seis peticiones están en la colección, en orden.
- ☑ Cada una responde lo que dice esta hoja.
- ☑ `go test ./...` en verde, con una prueba de integración.
- ☑ El código de la entidad principal está en capas y la estructura está explicada en `docs/decisiones.md`.
- ☑ Hay una captura de respaldo de cada paso.